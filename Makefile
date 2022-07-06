
NAME				:= exercise
NAMESPACE			:= default
HELM                := helm

help:
	@echo "install    : To install the exercise subchart"
	@echo "uninstall  : To uninstall the exercise subchart"

build:
	nerdctl build --namespace k8s.io  -t exercise:v1.1 .

## Installs Charts to kubernetes cluster
install:
	$(HELM) install --create-namespace -n $(NAMESPACE) $(NAME) ./chart/exercise

all: build install

rmi:
	nerdctl rmi -f --namespace k8s.io exercise:v1.1
## Uninstall charts
uninstall:
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)

all-remove: uninstall rmi


.PHONY : help install uninstall