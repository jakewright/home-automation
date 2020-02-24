// Code generated by jrpc. DO NOT EDIT.

package scenedef

import (
	context "context"

	rpc "github.com/jakewright/home-automation/libraries/go/rpc"
)

// Do performs the request
func (m *CreateSceneRequest) Do(ctx context.Context) (*CreateSceneResponse, error) {
	req := &rpc.Request{
		Method: "POST",
		URL:    "service.scene/scenes",
		Body:   m,
	}

	rsp := &CreateSceneResponse{}
	_, err := rpc.Do(ctx, req, rsp)
	return rsp, err
}

// Do performs the request
func (m *ReadSceneRequest) Do(ctx context.Context) (*ReadSceneResponse, error) {
	req := &rpc.Request{
		Method: "GET",
		URL:    "service.scene/scene",
		Body:   m,
	}

	rsp := &ReadSceneResponse{}
	_, err := rpc.Do(ctx, req, rsp)
	return rsp, err
}

// Do performs the request
func (m *ListScenesRequest) Do(ctx context.Context) (*ListScenesResponse, error) {
	req := &rpc.Request{
		Method: "GET",
		URL:    "service.scene/scenes",
		Body:   m,
	}

	rsp := &ListScenesResponse{}
	_, err := rpc.Do(ctx, req, rsp)
	return rsp, err
}

// Do performs the request
func (m *DeleteSceneRequest) Do(ctx context.Context) (*DeleteSceneResponse, error) {
	req := &rpc.Request{
		Method: "DELETE",
		URL:    "service.scene/scene",
		Body:   m,
	}

	rsp := &DeleteSceneResponse{}
	_, err := rpc.Do(ctx, req, rsp)
	return rsp, err
}

// Do performs the request
func (m *SetSceneRequest) Do(ctx context.Context) (*SetSceneResponse, error) {
	req := &rpc.Request{
		Method: "POST",
		URL:    "service.scene/scene/set",
		Body:   m,
	}

	rsp := &SetSceneResponse{}
	_, err := rpc.Do(ctx, req, rsp)
	return rsp, err
}