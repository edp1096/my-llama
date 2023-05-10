@echo off

del /s /q *.o
del /s /q *.a
del /s /q cgollama\*.o
del /s /q cgollama\*.a

cd llama.cpp
del /s /q *.o
rmdir build /s /q
rmdir openclblast /s /q >nul 2>&1
del /s /q clblast.zip >nul 2>&1
del /s /q opencl.zip >nul 2>&1
cd ..

del *.dll
del *.lib
del *.obj
del *.def
del *.exp