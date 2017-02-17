#!/usr/bin/env python
import shutil
import sys
import math
import datetime
import argparse
import tempfile
import os
import stat
from git import Repo

if sys.version_info[0] == 2:
    reload(sys)
    sys.setdefaultencoding('utf8')

BASE64_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
HEX_CHARS = "1234567890abcdefABCDEF"

def del_rw(action, name, exc):
    os.chmod(name, stat.S_IWRITE)
    os.remove(name)

def shannon_entropy(data, iterator):
    """
    Borrowed from http://blog.dkbza.org/2007/05/scanning-data-for-entropy-anomalies.html
    """
    if not data:
        return 0
    entropy = 0
    for x in (ord(c) for c in iterator):
        p_x = float(data.count(chr(x)))/len(data)
        if p_x > 0:
            entropy += - p_x*math.log(p_x, 2)
    return entropy


def get_strings_of_set(word, char_set, threshold=20):
    count = 0
    letters = ""
    strings = []
    for char in word:
        if char in char_set:
            letters += char
            count += 1
        else:
            if count > threshold:
                strings.append(letters)
            letters = ""
            count = 0
    if count > threshold:
        strings.append(letters)
    return strings

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def clone_repo(repo_path):
    project_path = tempfile.mkdtemp()
    Repo.clone_from(path, project_path)
    return project_path

def find_strings(project_path, outfile=None):

    repo = Repo(project_path)

    already_found = []
    already_searched = set()
    for remote_branch in repo.remotes.origin.fetch():
        branch_name = str(remote_branch).split('/')[1]
        try:
            repo.git.checkout(remote_branch, b=branch_name)
        except:
            pass

        prev_commit = None
        for curr_commit in repo.iter_commits():
            if not prev_commit:
                pass
            else:
                #avoid searching the same diffs
                hashes = str(prev_commit) + str(curr_commit)
                if hashes in already_searched:
                    prev_commit = curr_commit
                    continue
                already_searched.add(hashes)

                diff = prev_commit.diff(curr_commit, create_patch=True)
                for blob in diff:
                    #print i.a_blob.data_stream.read()
                    printableDiff = blob.diff.decode()
                    if printableDiff.startswith("Binary files"):
                        continue
                    foundSomething = False
                    lines = blob.diff.decode().split("\n")
                    for line in lines:
                        for word in line.split():
                            base64_strings = get_strings_of_set(word, BASE64_CHARS)
                            hex_strings = get_strings_of_set(word, HEX_CHARS)
                            for string in base64_strings:
                                b64Entropy = shannon_entropy(string, BASE64_CHARS)
                                if b64Entropy > 4.5:
                                    foundSomething = True
                                    printableDiff = printableDiff.replace(string, bcolors.WARNING + string + bcolors.ENDC)
                            for string in hex_strings:
                                hexEntropy = shannon_entropy(string, HEX_CHARS)
                                if hexEntropy > 3:
                                    foundSomething = True
                                    printableDiff = printableDiff.replace(string, bcolors.WARNING + string + bcolors.ENDC)
                    if foundSomething:
                        commit_time =  datetime.datetime.fromtimestamp(prev_commit.committed_date).strftime('%Y-%m-%d %H:%M:%S')
                        print(bcolors.OKGREEN + "Date: " + commit_time + bcolors.ENDC)
                        print(bcolors.OKGREEN + "Branch: " + branch_name + bcolors.ENDC)
                        print(bcolors.OKGREEN + "Commit: " + prev_commit.message + bcolors.ENDC)
                        print(printableDiff)

                        if outfile is not None and string not in already_found:

                            already_found.append(string)
                            outfile.write(string+'\n')

            prev_commit = curr_commit
    return project_path

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Find secrets hidden in the depths of git.')

    parser.add_argument('path', type=str, help='URL or PATH where to search for secrets')
    parser.add_argument('-o', '--outfile', nargs='?', type=argparse.FileType('w'), help='output file')

    args = parser.parse_args()

    if "http" in args.path:
        path = clone_repo(args.path)
    else:
        path = args.path

    project_path = find_strings(path, outfile=args.outfile)
    shutil.rmtree(project_path, onerror=del_rw)

