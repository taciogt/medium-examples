mypy:
	mypy -p domain
	#mypy -p domain -p api

test: mypy
	cd api; \
  	./manage.py test

#test: mypy
#	python -m unittest