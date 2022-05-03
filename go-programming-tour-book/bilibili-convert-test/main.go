package main

import (
	"fmt"
	"github.com/bilibili-convert-test/cmd"
	"log"
	"strings"
)

func main() {

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

	//testSetBegin()

	//urls := make([]string, 4)
	//urls[0] = "https://sjbz-fd.zol-img.com.cn/t_s320x510c5/g2/M00/05/0C/ChMlWl1BWGKIa5b1AAkDHph43SoAAMQfgALVicACQM2533.jpg"
	//urls[1] = "https://sjbz-fd.zol-img.com.cn/t_s320x510c5/g2/M00/05/0C/ChMlWl1BWGKIa5b1AAkDHph43SoAAMQfgALVicACQM2533.jpg"
	//urls[2] = "https://sjbz-fd.zol-img.com.cn/t_s320x510c5/g2/M00/05/0C/ChMlWl1BWGKIa5b1AAkDHph43SoAAMQfgALVicACQM2533.jpg"
	//urls[3] = "https://xiazai-fd.zol-img.com.cn/t_s960x600/g1/M01/03/06/Cg-4jVONmIiIa6NpAATdgesQtisAAN9YQLYqJcABN2Z899.jpg"
	//third_part.TestDownloadWithSingleThread("./", urls)

	log.Println("程序已退出")

	return
}

func testProgress() {
	f := float64(5612001) * 100 / 10000000
	cur := 5612001 * 100 / 10000000 / 2
	progress := fmt.Sprintf("%s%s%s", strings.Repeat("=", cur), ">", strings.Repeat(" ", 50-cur))
	fmt.Printf("当前下载进度: [%s] %.2f%%\n", progress, f)
}

func testSetBegin() {
	parts := make([]int, 10)
	for i := 0; i < 10; i++ {
		parts[i] = i
	}

	fmt.Println(parts[0:])
	fmt.Println(parts[1:])
	fmt.Println(parts[4:])
}
