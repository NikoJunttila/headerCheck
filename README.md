# headerCheck
<h1>How to use </h1> <br>
1. Clone repo <br>
2. Run go install or go build for executable <br>
3. In directory you want to check headers do command: headerCheck or ./headerCheck.exe if built (needs to be placed in same folder that running in so recommended to install) <br>
4.Then program checks all files for copyright headers and notifies if some need updating <br>
5. use flag -force to auto fix copyright headers. e.g: headerCheck -force<br>
Current flags <br>
 -force fixes files that are not correct<br>
 --ignore="" Specify folders/files to ignore -ignore="vendor" -ignore="node_modules" you can use this multiple times to ignore many folders/files. Filepath is relational to directory where you start program so -ignore="tests\test.go" would ignore file inside tests directory <br>
--author="Niko Junttila niko.junttila2@centria.fi" adds default author name/email if for some reason project is not in git/merc repo<br>
--year="2022-2023" adds default year <br>
-suffix=".js,.py,.cpp" to only go through those files and skip others <br>
--newSuf=".haskell?" add new suffix to default list if not already included. <br>
-verbose  Prints more errors/info about whats happening <br>
--noHeader ignores files that don't have a header already if you want to only check already existing headers. <br>
-single="string" only checks this file and ignores rest <br>
<br>
-usage flag for list of flags <bold>AND </bold> location of global exe for possible template.txt
<br>
For custom template modify template.txt. You can have multiple.<br>
1st. checks flag -template="src" for custom template <br>
2nd. checks current working directory where you call the executable So you can place project specific template there <br>
3rd. checks global executable folder for template.txt (windows default location is $user/go/bin/ check with -usage if unsure)<br>
4th. uses hardcoded default Centria template if nothing else was found.<br>
*checks for .hg file if not found defaults to git* <br> 
currently checks these suffixes: .go, .cpp, .c, .h, .hpp, .js, .ts, .cs, .java, .rs, .qml, .css, .qss, .scala, .kt, .jsx, .tsx, .swift, .zig, .py, .exs, .hmtl, .rb, .lua, .ml, mli
<br>
<h2>for linux users</h2>
<br>
1. go build <br>
2. sudo mv executable_name /usr/local/bin/ <br>
3. headerCheck -force <br>
