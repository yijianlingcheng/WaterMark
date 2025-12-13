package pkg

import (
	"bufio"
	"os"
	"strings"
)

type CSV struct {
	name    string
	folder  string
	path    string
	headers []string
	data    [][]string
	hasBOM  bool
}

// 创建csv.
func CreateCSV(name, folder string, hasBOM bool) *CSV {
	return &CSV{
		hasBOM: hasBOM,
		name:   name,
		folder: folder,
		path:   folder + name,
	}
}

func (c *CSV) getLine() string {
	if IsWindows() {
		return "\r\n"
	}

	return "\n"
}

// 设置csv头.
func (c *CSV) SetHeaders(headers []string) {
	c.headers = headers
}

// 添加数据.
func (c *CSV) AddData(data [][]string) {
	c.data = data
}

// 生成文件.
func (c *CSV) Generate() EError {
	file, err := os.Create(c.path)
	if err != nil {
		errmsg := c.path + ":创建csv文件失败:" + err.Error()
		// 返回失败
		return NewErrors(CSV_CREATE_ERROR, errmsg)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// 写入BOM头
	if c.hasBOM {
		_, err = writer.WriteString("\xEF\xBB\xBF")
		if err != nil {
			return NewErrors(CSV_WRITE_HEADER_ERROR, "csv写入bom头失败")
		}
	}

	if len(c.headers) > 0 {
		// 写入数据
		_, err = writer.WriteString(strings.Join(c.headers, ",") + c.getLine())
		if err != nil {
			errmsg := c.path + ":写入headers失败:" + strings.Join(c.headers, ",") + ",原因:" + err.Error()
			// 返回失败
			return NewErrors(CSV_WRITE_HEADER_ERROR, errmsg)
		}
	}

	for _, item := range c.data {
		_, err = writer.WriteString(strings.Join(item, ",") + c.getLine())
		if err != nil {
			errmsg := c.path + ":写入data失败:" + strings.Join(item, ",") + ",原因:" + err.Error()
			// 返回失败
			return NewErrors(CSV_WRITE_DATA_ERROR, errmsg)
		}
	}
	writer.Flush()

	return NoError
}
