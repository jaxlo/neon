dev:
	cc src/main.c -lraylib -lGL -lm -lpthread -ldl -lrt -lX11
	./a.out
	rm -f a.out
