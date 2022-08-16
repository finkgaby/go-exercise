
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



rmi:
	nerdctl rmi -f --namespace k8s.io exercise:v1.1
## Uninstall charts
uninstall:
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)

all-remove: uninstall rmi

all: uninstall build install

.PHONY : help install uninstall