cd binding

cl.exe /nologo /LD binding.cpp /Fe:./

gendef ./binding.dll
dlltool -k -d ./binding.def -l ./libbinding.a

mv ./binding.dll ../binding.dll
mv ./libbinding.a ../libbinding.a

rm -f binding.obj
rm -f binding.exp
rm -f binding.lib
rm -f binding.def
