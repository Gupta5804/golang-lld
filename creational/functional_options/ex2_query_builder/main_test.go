package main

import (
	"reflect"
	"testing"
)

func TestNewSelectQuery_Defaults(t *testing.T){
	query := NewSelectQuery("users")
	expectedCols  := []string{"*"}
	if !reflect.DeepEqual(query.columns,expectedCols){
		t.Errorf("Expected columns :'*', instead got %v", query.columns)
	}
	if query.where != ""{
		t.Errorf("Expected where : empty, got %s", query.where)
	}
	if query.orderBy != ""{
		t.Errorf("Expected orderBy : empty, instead got %s", query.orderBy)
	}
	if query.limit != 0 {
		t.Errorf("Expected LIMIT:0 , instead got %d", query.limit)
	}
}
func TestNewSelectQuery_WithOptions(t *testing.T){
	query := NewSelectQuery("orders",
		WithLimit(5),
		WithWhere("status = 'paid'"),
	)
	expectedCols  := []string{"*"}
	if !reflect.DeepEqual(query.columns,expectedCols){
		t.Errorf("Expected columns :'*', instead got %v", query.columns)
	}
	if query.where != "status = 'paid'"{
		t.Errorf("Expected where : status = 'paid', got %s", query.where)
	}
	if query.orderBy != ""{
		t.Errorf("Expected orderBy : empty, instead got %s", query.orderBy)
	}
	if query.limit != 5{
		t.Errorf("Expected LIMIT:5 , instead got %d", query.limit)
	}
}