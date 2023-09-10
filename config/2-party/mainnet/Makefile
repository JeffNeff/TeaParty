stop:
	@kubectl delete -f 0-teabarrel.yaml && kubectl delete -f 3-adams.yaml 
start:
	@kubectl apply -f 0-teabarrel.yaml && kubectl apply -f 3-adams.yaml
watch:
	@kubectl logs -l  app=adams -f