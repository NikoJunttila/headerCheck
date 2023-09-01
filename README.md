# headerCheck
<h1>How to use </h1> <br>
1. Clone repo <br>
2. Run go install or go build for executable <br>
3. In directory you want to check headers do command: headerCheck or ./headerCheck.exe if built (needs to be placed in same folder that running in so recommended to install) <br>
4.Then program checks all files for copyright headers and notifies if some need updating <br>
5. use flag -force to auto fix copyright headers. e.g: headerCheck -force<br>
Current flags -force,--ignore="folder or file" if you want to ignore a folder/file from check, <br>
--author="if no git history add default name" --year="add default year i.e "2022-2023"" <br>
-suffix=".js,.py,.cpp" to only go through those files and skip others
--newSuf=".js" add new suffix to default list if not already included.
<br>
<h2>for linux users</h2>
<br>
1. go build <br>
2. sudo mv executable_name /usr/local/bin/ <br>
3. headerCheck -force <br>
<br>
-usage flag for list of flags <bold>AND </bold> location of global exe for possible template.txt
<br>
For custom template modify template.txt. You can have multiple.<br>
1st. checks flag -template="src" for custom template <br>
2nd. checks current working directory where you call the executable So you can place project specific template there <br>
3rd. checks global executable folder for template.txt (windows default location is $user/go/bin/ check with -usage if unsure)<br>
4th. uses hardcoded default Centria template if nothing else was found.<br>
*checks for .hg file if not found defaults to git*
