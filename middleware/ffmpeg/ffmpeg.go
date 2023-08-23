package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// 抽取指定路径下的视频画面帧作为视频封面图，保存到指定路径下
func ExtractThumbnail(videoPath, thumbnailRoot string) error {
	// 检查图片文件是否存在，不存在即创建
	// ffmpeg不支持写入结果到原本不存在的文件
	CreateCoverFile(thumbnailRoot)

	// 抽取帧的时间点,第一帧
	frameTime := "00:00:00.0"

	// 执行 ffmpeg 命令
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", frameTime, "-vframes", "1", thumbnailRoot, "-y")

	fmt.Println("save cover cmd ", cmd)

	_, err := cmd.Output()

	if err != nil {
		return err
	}

	return nil
}

func CreateCoverFile(thumbnailRoot string) error {
	// 获取父级目录
	dir := filepath.Dir(thumbnailRoot)

	err := os.MkdirAll(dir, os.ModePerm) // 递归创建缺失的父级目录
	if err != nil {
		fmt.Println("创建父级目录失败:", err)
		return err
	}

	_, err = os.Stat(thumbnailRoot)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(thumbnailRoot)
		if err != nil {
			fmt.Println("创建文件失败:", err)
			return err
		}
		defer file.Close()
		fmt.Println("文件已创建")
	} else if err != nil {
		// 其他错误
		fmt.Println("发生错误:", err)
		return err
	} else {
		// 文件已存在
		fmt.Println("文件已存在")
	}
	return err
}
