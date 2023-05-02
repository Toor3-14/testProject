package model
 
import (
	"errors"
)
 
var ErrNoRecord = errors.New("models: not found")
 
type Users struct {
	Name   string
	Age int
}