package bayesian

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func isTheSameName(name1 string, name2 string) bool {
	n1 := []byte(name1)
	n2 := []byte(name2)

	if len(n1) != len(n2) {
		return false
	}

	for i, v := range n1 {
		if v != n2[i] {
			return false
		}
	}

	return true
}

func BayesianClassifyFromFile(fileName string) {
	fd, err := os.Open(fileName)
	defer fd.Close()
	if err != nil {
		fmt.Println("Read", fileName,  "error:", err)
		return
	}
	buff := bufio.NewReader(fd)

	// 读取属性名称列
	firstLine, _, eof := buff.ReadLine()
	if eof == io.EOF {
		fmt.Println("Empty file!")
		return
	}

	attrs := strings.Split(string(firstLine), ",")
	//fmt.Println("Attribute names: ", attrs[1:len(attrs)-1])
	//fmt.Println("Result name: ", attrs[len(attrs)-1])

	resultValues := []string{"是", "否"}

	// 创建新的对象 结果列的名称，结果列可以取得值，属性列的名称
	bayesian := NewBayesian(attrs[len(attrs)-1], resultValues, attrs[1:len(attrs)-1])

	// 添加所有的属性值
	bayesian.AddAttribute(0, "色泽", []string{"青绿", "乌黑", "浅白"})
	bayesian.AddAttribute(1, "根蒂", []string{"蜷缩", "稍蜷", "硬挺"})
	bayesian.AddAttribute(2, "敲声", []string{"浊响", "沉闷", "清脆"})
	bayesian.AddAttribute(3, "纹理", []string{"清晰", "稍糊", "模糊"})
	bayesian.AddAttribute(4, "脐部", []string{"凹陷", "稍凹", "平坦"})
	bayesian.AddAttribute(5, "触感", []string{"硬滑", "软粘"})

	// 设置想要计算的属性值
	bayesian.TaskInput([]string{"青绿", "稍蜷", "浊响", "清晰", "凹陷", "硬滑"})

	// 训练模型
	for {
		id := 0
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}

		values := strings.Split(string(data), ",")
		//fmt.Println(id, "\t", "Attribute names: ", values[1:len(values)-1], "Result name: ", values[len(values)-1])
		entity := &Entity{id, values[1:len(values)-1],values[len(values)-1]}
		bayesian.Exercise(entity)
		id++
	}

	// 输出结果
	out := bayesian.TaskOutput()
	fmt.Printf("Attributes: %s, result: %s\n", bayesian.Task.Values, out)
}
