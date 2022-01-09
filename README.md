# NAME

sttemp - simple template manager

# SYNOPSIS

**sttemp** \[**-h**\]

**sttemp** \[**-e**\] **template** \[**output**\]

# DESCRIPTION

**sttemp** copy file from your template directory and fill all
substitutions with user input, or environment variables. Template
directory and substitution pattern described in **src/config.h** file,
in suckless style. You can use **\~** as **first** symbol of your path,
for relative path from your home directory.

# OPTIONS

**-h**

:   displays help information

**-e**

:   use environment variables for fill substitution in template, if
    **sttemp** can\'t find environment variable with needed name, it
    will use user input

**template**

:   name of file in your templates directory, which will use to create
    new file

**output**

:   name of target file, if omit - sttemp will use template name

# AUTHOR

Anton Konyahin \<me\@konyahin.xyz>
