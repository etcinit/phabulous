package mocks

//go:generate mockgen -destination=generated.go -package=mocks github.com/etcinit/phabulous/app/interfaces Bot,Message,HandlerTuple,Module,Connector
