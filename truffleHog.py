import shutil, sys, math, string, datetime, argparse, tempfile
from git import Repo

# needed to get list of orgs using github api
import requests

if sys.version_info[0] == 2:
    reload(sys)  
    sys.setdefaultencoding('utf8')

BASE64_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
HEX_CHARS = "1234567890abcdefABCDEF"

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
            if count > 20:
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

def find_strings(git_url):
    project_path = tempfile.mkdtemp()

    Repo.clone_from(git_url, project_path)

    repo = Repo(project_path)


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
                diff = prev_commit.diff(curr_commit, create_patch=True)
                for blob in diff:
                    #print i.a_blob.data_stream.read()
                    printableDiff = blob.diff.decode() 
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
                    
            prev_commit = curr_commit
    shutil.rmtree(project_path)

def get_org_repos(orgname):
    response = requests.get(url='https://api.github.com/users/' + orgname + '/repos')
    json = response.json()
    for item in json:
        if item['private'] == False:
            print('searching ' + item["html_url"])
            find_strings(item["html_url"])

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Find secrets hidden in the depths of git.')
    subparsers = parser.add_subparsers(help='commands')
    # Parse org repo search command
    org_parser = subparsers.add_parser('org', help='Search inside all repositories for an organization (github.com/uber, github.com/google, etc)')
    org_parser.add_argument('orgname', type=str, action='store', help='Github Organization name (uber, yubico, etc)')
    # Parse git repo command
    git_parser = subparsers.add_parser('git', help='Search inside a singpe repository')
    git_parser.add_argument('git_url', type=str, action='store', help='URL for secret searching')
    args = parser.parse_args()
    if hasattr(args, 'orgname'):
        print('searching for repos in the github organization: ' + args.orgname)
        get_org_repos(args.orgname)
    if hasattr(args, 'git_url'):
        find_strings(args.git_url)
