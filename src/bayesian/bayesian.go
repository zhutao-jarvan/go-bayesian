package bayesian

import (
	"errors"
	"fmt"
)

type Probability struct {
	numerator uint64 // 分子
	denominator uint64 // 分母
}

type ResultProbStatistics struct {
	ResultValue string
	Total Probability
	Attributes []Probability
	P float32
}

type Attribute struct {
	Name string 		// 属性
	Values []string     // 取值列表
}

type Entity struct {
	Id int
	Values []string
	Result string
}

type BayesianTask struct {
	Values []string
	ResultProbs []ResultProbStatistics  // 数量等于Result的可取值数量
}

type Bayesian struct {
	Result Attribute
	Attributes []Attribute

	Task BayesianTask
}

// 创建一个新的对象，准备进行计算
// @result 结果列的名称
// @resultValues 结果列的取值范围
// @attributeNames 所有属性列的名称
func NewBayesian(result string, resultValues []string, attributeNames []string) *Bayesian {
	bayesian := &Bayesian{
		Attributes: make([]Attribute, 0),
	}

	bayesian.Result.Name = result
	bayesian.Result.Values = resultValues

	for _, attributeName := range attributeNames {
		attr := Attribute{attributeName, make([]string, 0)}
		bayesian.Attributes = append(bayesian.Attributes, attr)
	}

	return bayesian
}

// 添加一个属性列的取值范围
// @attrID 属性是第几列
// @attribute 属性名称
// @values 改属性列的取值范围
func (bayesian *Bayesian) AddAttribute(attrID uint32, attribute string, values []string) error {
	if len(bayesian.Attributes) <= int(attrID) {
		return errors.New(fmt.Sprintf("Attribute id %d is too great, max: %d\n", int(attrID), len(bayesian.Attributes)))
	}

	// TODO check attribute name
	for _, value := range values {
		bayesian.Attributes[attrID].Values = append(bayesian.Attributes[attrID].Values, value)
	}

	return nil
}

// 输入一行属性值进行统计
func (bayesian *Bayesian) Exercise(entity *Entity) error {
	task := &bayesian.Task

	for i, _ := range task.ResultProbs {
		p := &task.ResultProbs[i]
		if isTheSameName(p.ResultValue, entity.Result) {
			// 结果一致
			// 总的统计概率各+1
			p.Total.denominator++
			p.Total.numerator++

			// 其他各属性分母+1, 分子相等者+1
			for j, value := range entity.Values {
				p.Attributes[j].denominator++

				if isTheSameName(task.Values[j], value) {
					p.Attributes[j].numerator++
				}
			}
		} else {
			// 结果不一致
			// 总的统计概率分母+1
			p.Total.denominator++

			// 其他各属性不变
		}
	}
	fmt.Println(task.ResultProbs)

	return nil
}

// 根据一行属性值来新建一个任务，需要初始化计算任务统计技术的对象
// @values 待计算的某一个属性值的行
func (bayesian *Bayesian) TaskInput(values []string) error {
	// 存储待计算的目标属性
	if len(values) != len(bayesian.Attributes) {
		return errors.New(fmt.Sprintf("Task Values number %d is error, need: %d\n", len(values), len(bayesian.Attributes)))
	}

	bayesian.Task.Values = values

	// 准备需要计算概率的存储空间
	task := &bayesian.Task
	task.ResultProbs = make([]ResultProbStatistics, len(bayesian.Result.Values))
	for i, value := range bayesian.Result.Values {
		task.ResultProbs[i].ResultValue = value

		// 每一种结果的情况都需要各个属性的列
		task.ResultProbs[i].Attributes = make([]Probability, len(bayesian.Attributes))
		for j, _ := range bayesian.Attributes {
			task.ResultProbs[i].Attributes[j] = Probability{}
		}
	}

	return nil
}

// 计算出概率最高的记过
func (bayesian *Bayesian) TaskOutput() string {
	var resultIdx int
	task := bayesian.Task
	PResult := float32(-1)

	for i, p := range task.ResultProbs {
		numerator := p.Total.numerator // 分子
		denominator := p.Total.denominator // 分母

		for j, _ := range task.ResultProbs[i].Attributes {
			numerator *= p.Attributes[j].numerator
			denominator *= p.Attributes[j].denominator // overflow check
		}

		p.P = float32(numerator) / float32(denominator)
		if p.P > PResult {
			PResult = p.P
			resultIdx = i
		}

		//fmt.Println(numerator, denominator)
		fmt.Printf("Result %s probbility [%d/%d] %0.08f\n", task.ResultProbs[i].ResultValue, numerator, denominator, p.P)
	}

	return task.ResultProbs[resultIdx].ResultValue
}




