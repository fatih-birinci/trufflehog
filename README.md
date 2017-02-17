# Truffle Hog
Searches through git repositories for high entropy strings, digging deep into commit history and branches. This is effective at finding secrets accidentally committed that contain high entropy.

```
python truffleHog.py https://github.com/dxa4481/truffleHog.git
```

![Example](https://i.imgur.com/YAXndLD.png)

## Setup
The only requirement is GitPython, which can be installed with the following
```
pip install -r requirements.txt
```

## How it works

This module will go through the entire commit history of each branch, and check each diff from each commit, and evaluate the shannon entropy for both the base64 char set and hexidecimal char set for every blob of text greater than 20 characters comprised of those character sets in each diff. If at any point a high entropy string >20 characters is detected, it will print to the screen. 

## Run

It is possible to run it in a remote git repository or with a local git repository:

`
python truffleHog.py source_folder
`

or

`
python truffleHog.py http://...
`

## Output

This module will output print in the diffs where the key is found highlighting the found string.
It is possible to make it ouput a file with the argument:

	`-o` or `-outfile`

This file will contain all unique strings found. One per line.

## Wishlist

- ~~A way to detect and not scan binary diffs~~
- ~~Don't rescan diffs if already looked at in another branch~~
- ~~Output a file with all found keys~~