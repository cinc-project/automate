package pmt

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"regexp/syntax"
	"strconv"

	"github.com/manifoldco/promptui"
)

type Prompt interface {
	Confirm(question, trueOption, falseOption string) (result bool, err error)
	Select(question string, options ...string) (result string, err error)
	InputInt64(label string) (resultVal int64, err error)
	InputInt(label string) (resultVal int, err error)
	InputFloat64(label string) (resultVal float64, err error)
	InputStringMinMax(label string, minlen int, maxlen int) (result string, err error)
	InputWordDefault(label string, defaultVal string) (result string, err error)
	InputWord(label string) (result string, err error)
	InputStringRegx(label string, regexCheck string) (result string, err error)
	InputString(label string) (result string, err error)
	InputStringRegxDefault(label string, regexCheck string, defaultVal string) (result string, err error)
	InputStringDefault(label string, defaultVal string) (result string, err error)
}

type PromptImp struct {
	Stdin  io.ReadCloser
	Stdout io.WriteCloser
}

func PromptFactory(in io.ReadCloser, out io.WriteCloser) *PromptImp {
	return &PromptImp{
		Stdin:  in,
		Stdout: out,
	}
}

func (p *PromptImp) Confirm(question, trueOption, falseOption string) (resultVal bool, err error) {
	prompt := promptui.Select{
		Label:  question,
		Items:  []string{trueOption, falseOption},
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
	}

	_, result, err := prompt.Run()
	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}

	switch result {
	case trueOption:
		return true, nil
	case falseOption:
		return false, nil
	default:
		return false, nil
	}
}

func (p *PromptImp) Select(question string, options ...string) (result string, err error) {
	prompt := promptui.Select{
		Label:  question,
		Items:  options,
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
	}

	_, result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) SelectAdd(question string, addLabel string, options ...string) (result string, err error) {
	index := -1

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    question,
			Items:    options,
			AddLabel: addLabel,
		}

		index, result, err = prompt.Run()
		if index == -1 {
			options = append(options, result)
		}
	}
	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputInt64(label string) (resultVal int64, err error) {
	validate := func(input string) error {
		_, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("invalid int64")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err := prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}

	resultVal, err = strconv.ParseInt(result, 10, 64)
	if err != nil {
		err = fmt.Errorf("parse int64 failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputInt(label string) (resultVal int, err error) {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("invalid int")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err := prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}

	resultVal, err = strconv.Atoi(result)
	if err != nil {
		err = fmt.Errorf("parse int failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputFloat64(label string) (resultVal float64, err error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid float64")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err := prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}

	resultVal, err = strconv.ParseFloat(result, 64)
	if err != nil {
		err = fmt.Errorf("parse float64 failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputWord(label string) (result string, err error) {
	return p.InputWordDefault(label, "")
}

func (p *PromptImp) InputWordDefault(label string, defaultVal string) (result string, err error) {
	validate := func(input string) error {
		isWord := true
		for _, v := range input {
			isWord = syntax.IsWordChar(v)
			if !isWord {
				return fmt.Errorf("invalid word %v", err)
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  defaultVal,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputStringMinMax(label string, minlen int, maxlen int) (result string, err error) {
	validate := func(input string) error {
		if len(input) < minlen {
			return fmt.Errorf("smaller than %v", minlen)
		} else if len(input) > maxlen {
			return fmt.Errorf("larger than %v", maxlen)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputStringRegx(label string, regexCheck string) (result string, err error) {
	validate := func(input string) error {
		var check = regexp.MustCompile(regexCheck).MatchString
		if !check(input) {
			return fmt.Errorf("regex check failed (%v) ", regexCheck)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputStringRegxDefault(label string, regexCheck string, defaultVal string) (result string, err error) {
	validate := func(input string) error {
		var check = regexp.MustCompile(regexCheck).MatchString
		if !check(input) {
			return fmt.Errorf("regex check failed (%v) ", regexCheck)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  defaultVal,
		Stdin:    p.Stdin,
		Stdout:   p.Stdout,
	}

	result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputString(label string) (result string, err error) {
	prompt := promptui.Prompt{
		Label:  label,
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
	}

	result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}

func (p *PromptImp) InputStringDefault(label string, defaultVal string) (result string, err error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
		Stdin:   p.Stdin,
		Stdout:  p.Stdout,
	}

	result, err = prompt.Run()

	if err != nil {
		err = fmt.Errorf("prompt failed %v", err)
		return
	}
	return
}
