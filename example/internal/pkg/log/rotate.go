package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// RotateConfig 定义日志文件轮转的配置
type RotateConfig struct {
	MaxSize    int  // 单个日志文件的最大大小（MB）
	MaxAge     int  // 日志文件的最大保留天数
	MaxBackups int  // 最大保留的日志文件数
	Compress   bool // 是否压缩旧的日志文件
}

// DefaultRotateConfig 返回默认的日志轮转配置
func DefaultRotateConfig() *RotateConfig {
	return &RotateConfig{
		MaxSize:    100,   // 默认每个日志文件最大100MB
		MaxAge:     7,     // 默认保留7天
		MaxBackups: 10,    // 默认保留10个备份
		Compress:   false, // 默认不压缩
	}
}

// rotateManager 管理日志文件的轮转
type rotateManager struct {
	lock   sync.Mutex
	config *RotateConfig
	sizes  map[string]int64 // 跟踪每个日志文件的大小
}

var (
	rotateManagerInstance *rotateManager
	rotateOnce            sync.Once
)

// GetRotateManager 获取日志轮转管理器的单例实例
func GetRotateManager() *rotateManager {
	rotateOnce.Do(func() {
		rotateManagerInstance = &rotateManager{
			config: DefaultRotateConfig(),
			sizes:  make(map[string]int64),
		}
	})
	return rotateManagerInstance
}

// Configure 配置日志轮转管理器
func (r *rotateManager) Configure(config *RotateConfig) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if config != nil {
		r.config = config
	}
}

// Write 写入内容到日志文件，如果需要则进行轮转
func (r *rotateManager) Write(path string, data []byte) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 检查文件是否需要轮转
	if r.shouldRotate(path) {
		if err := r.rotate(path); err != nil {
			return fmt.Errorf("rotate log file error: %w", err)
		}
	}

	// 打开文件（追加模式）
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open log file error: %w", err)
	}
	defer f.Close()

	// 写入数据
	n, err := f.Write(data)
	if err != nil {
		return fmt.Errorf("write log data error: %w", err)
	}

	// 更新文件大小记录
	r.sizes[path] += int64(n)

	return nil
}

// shouldRotate 判断是否需要轮转日志文件
func (r *rotateManager) shouldRotate(path string) bool {
	// 获取当前文件大小
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，不需要轮转
			r.sizes[path] = 0
			return false
		}
		// 出错时保守处理，不进行轮转
		fmt.Fprintf(os.Stderr, "Failed to stat log file %s: %v\n", path, err)
		return false
	}

	// 更新文件大小记录
	r.sizes[path] = info.Size()

	// 检查文件大小是否超过限制
	return r.sizes[path] >= int64(r.config.MaxSize*1024*1024)
}

// rotate 执行日志文件轮转
func (r *rotateManager) rotate(path string) error {
	// 如果文件不存在，不需要轮转
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// 生成轮转后的文件名
	backupName := r.getBackupFilename(path)

	// 关闭可能打开的文件
	// 重命名文件
	if err := os.Rename(path, backupName); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// 重置文件大小记录
	r.sizes[path] = 0

	// 执行清理
	if err := r.cleanup(path); err != nil {
		// 仅记录错误，不中断流程
		fmt.Fprintf(os.Stderr, "Failed to cleanup log files: %v\n", err)
	}

	return nil
}

// getBackupFilename 生成轮转后的备份文件名
func (r *rotateManager) getBackupFilename(path string) string {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	timestamp := time.Now().Format("20060102-150405")
	return filepath.Join(dir, fmt.Sprintf("%s.%s", filename, timestamp))
}

// cleanup 清理过期的日志文件
func (r *rotateManager) cleanup(path string) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// 读取目录内容
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	var backups []string
	// 找出属于该日志文件的所有备份
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), base+".") {
			backups = append(backups, filepath.Join(dir, file.Name()))
		}
	}

	// 按照修改时间排序
	sort.Slice(backups, func(i, j int) bool {
		iInfo, _ := os.Stat(backups[i])
		jInfo, _ := os.Stat(backups[j])
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	// 删除超过最大备份数量的文件
	if len(backups) > r.config.MaxBackups && r.config.MaxBackups > 0 {
		for i := r.config.MaxBackups; i < len(backups); i++ {
			if err := os.Remove(backups[i]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to remove old log file %s: %v\n", backups[i], err)
			}
		}
	}

	// 如果设置了最大保留天数，删除过期文件
	if r.config.MaxAge > 0 {
		cutoff := time.Now().Add(-time.Duration(r.config.MaxAge) * 24 * time.Hour)

		for _, file := range backups {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoff) {
				if err := os.Remove(file); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove expired log file %s: %v\n", file, err)
				}
			}
		}
	}

	return nil
}
