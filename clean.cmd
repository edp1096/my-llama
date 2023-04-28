@echo off

del /s /q *.o
del /s /q *.a
del /s /q cgollama\*.o
del /s /q cgollama\*.a

cd llama.cpp
del /s /q *.o
rmdir build /s /q
cd ..

del *.dll
del *.lib
del *.obj
del *.def
del *.exp