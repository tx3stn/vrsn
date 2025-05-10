Output demo.gif
Require git
Set Theme "OneDark"
Set Height 1000

Hide
Type "alias vrsn=$(pwd)/vrsn"
Enter
Type "clear"
Enter
Show

Set TypingSpeed 0.2s
Type "vrsn"
Enter 2
Sleep 2s

Type "vrsn check"
Enter 2
Sleep 3s

Type "vrsn bump"
Enter
Sleep 1s
Down 1
Sleep 0.5s
Down 1
Sleep 0.5s
Up 1
Sleep 0.5s
Up 1
Sleep 2.5s
Enter 2

Type "vrsn check"
Enter 2
Sleep 3s

Set TypingSpeed 0s
Hide
Type "git checkout ./VERSION"
Enter
