package main

import "fmt"

type Document struct {
	currentState State
}

func NewDocument() *Document {
	doc := &Document{}
	doc.SetState(draftState)
	return doc
}
func (d *Document) SetState(s State) {
	d.currentState = s
}
func (d *Document) GetStateName() string{
	return d.currentState.Name()
}
func (d *Document) RequestReview() error{
	return d.currentState.RequestReview(d)
}
func (d *Document) Approve() error{
	return d.currentState.Approve(d)
}

type State interface {
	Name() string
	RequestReview(d *Document) error
	Approve(d *Document) error
}

var draftState = &DraftState{}
var reviewState = &ReviewState{}
var publishedState = &PublishedState{}

// Draft State
type DraftState struct{}

func (s *DraftState) Name() string {
	return "Draft"
}
func (s *DraftState) RequestReview(d *Document) error {
	d.SetState(reviewState)
	return nil
}
func (s *DraftState) Approve(d *Document) error {
	return fmt.Errorf("cannot approve a draft directly")
}

// Review State
type ReviewState struct{}

func (s *ReviewState) Name() string {
	return "Review"
}
func (s *ReviewState) RequestReview(d *Document) error {
	return fmt.Errorf("already in review")
}
func (s *ReviewState) Approve(d *Document) error {
	d.SetState(publishedState)
	return nil
}

// Published State
type PublishedState struct{}
func (s *PublishedState) Name() string {
	return "Published"
}
func (s *PublishedState) RequestReview(d *Document) error {
	return fmt.Errorf("already published")
}
func (s *PublishedState) Approve(d *Document) error {
	return fmt.Errorf("already published")
}