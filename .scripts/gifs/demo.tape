Set Height 1000
Source '/vhs/.scripts/gifs/base.tape'

Output demo.gif

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
Sleep 5s

Set TypingSpeed 0s
Hide
Type "git checkout ./VERSION"
Enter
