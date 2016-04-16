.PHONY: clean

TARGET=faceglass_server

$(TARGET): libfoo.a
	go build .

libfoo.a: foo.o cfoo.o
	ar r $@ $^

%.o: %.cpp
	g++ -O2 -o $@ -c $^

clean:
	rm -f *.o *.so *.a $(TARGET)
