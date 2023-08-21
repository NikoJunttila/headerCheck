# headerCheck
<h1>How to use </h1> <br>
1. Clone repo <br>
2. Run go install or go build for executable <br>
3. In directory you want to check headers do command: headerCheck or ./headerCheck.exe if built (needs to be placed in same folder that running in so recommended to install) <br>
4.Then program checks all files for copyright headers and notifies if some need updating <br>
5. use flag -force to auto fix copyright headers. e.g: headerCheck -force<br>
Current flags -force,--ignore="folder" if you want to ignore a folder from check, <br>
--author="if no git history add default name" --year="add default year i.e "2022-2023"" <br>
<br>
for linux users<br>
1. go build <br>
2. sudo mv executable_name /usr/local/bin/ <br>
3. headerCheck -force <br>
*checks for .git file if not found defaults to hg*
