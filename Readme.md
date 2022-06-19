<picture>
  <source media="(prefers-color-scheme: dark)" srcset="images/logo-dark.png">
  <img alt="Different logo for light and dark themes" src="images/logo.png">
</picture>

**claro** is a [GitHub Classroom](https://classroom.github.com) CLI for teachers

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

# Overview

**claro** (**cla**ss**ro**om) is a cli tool that offers a simple interface that allows the teacher to clone all student repositories at once for grading and then send grades at once to all these repositories.

**claro** was inspired by the [Git-Gud-tool](https://github.com/NikolaiMagnussen/Git-Gud-tool) and it was created to make my life simpler as a teacher who relies on Github Classroom. It is suitable for scenarios where the teacher needs to manually grade each assignment and the Github Classroom autograding is not an option. 

**claro** relies on [Github's REST API](https://docs.github.com/en/rest) because currently Github Classroom does not offer an API and it has features that are not present in [Github Classroom Assistant](https://classroom.github.com/assistant). The best one: it is a [CLI](https://clig.dev) :heart:!


**claro** provides:
- mass clone of Github repositories
- a template in Markdown for grading (it creates one file per repository)
- customization (commit message, grade sheet title, grading file, etc.)
- access to [Github Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) from operating system keyring (i.e. macOS Keychain, Gnome Keyring), environment var or claro's config file

## Requirements
- [git](https://git-scm.com/docs/git) command line application configured properly
   ```bash
   git config --global user.email "you@example.com"
   git config --global user.name "Your Name"
   ```
- [Git credential store](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage) configured properly
  - See [section below](#authenticating-with-git-command-line-to-access-repositories-on-github)

## Download

- [Download Binaries for different platforms](https://github.com/emersonmello/claro/releases/latest)
  - macOS
    - darwin-amd64
    - darwin-arm64
  - Linux
    - linux-386
    - linux-amd64
    - linux-arm64
  

## Usage

### Straight workflow

1. Create a Github Classroom assignment with a distinguish repository prefix 
   - Example: `2022-01-assignment-01`
2. Clone all repositories from a Github Classroom organization with a specific assignment prefix

    ![cloning](images/clone.gif)

3. Grade each student's work and write down the feedback in the respective Markdown file

    ![grading](images/grading.gif)

4. Push the grading

    ![pushing](images/push.gif)

## Other commands

![claro help message](images/claro.png)

### Pull students repositories

**claro** offers a `pull` command suitable for workflows where students make incremental deliveries to the same GitHub Classroom repository. 

- Example:
  ```bash
  claro pull 2022-01-assignment-01
  ```
### List all repositories which start with a specific prefix

Before cloning multiple repositories, you might want to check which repositories will be cloned with a specific prefix. **claro** offers the `list` command for that.

- Example:
  ```bash
  claro list 2022-01-assignment-01
  ```

### Customize commit message, grading filename, grading string

The **claro** default strings are:

- **Grading filename:** `GRADING.md`
  -  It will be created and pushed to student repository
- **Grade string:** `Grade: `
  - It will be inside grading file 
- **Commit message:** `Graded project, the file containing the grade is in the root directory`
  - It is the commit message
- **Grade sheet title** `Feedback`
  - It will be inside grading file as title 1 (# Feedback)

You can change the values using **claro** `config` command. The new values will be stored in the **claro**'s config file (default `$HOME/.claro.env`).
- Examples:
  - `claro config filename "Correcao.md"`
  - `claro config grade "Nota: "`
  - `claro config message "Correção finalizada, veja arquivo na raiz do repositório"`
  - `claro config title "Comentários"`

## Storing claro's GitHub Personal Access Token in the OS keyring

**claro** only supports HTTPS remote URL (git over SSH is so [annoying](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#authenticating-with-the-command-line)). 

**claro** consumes [Github's REST API](https://docs.github.com/en/rest) to fetch the list of repositories (assignments) from a GitHub Classroom organization. So, **claro** uses a [GitHub Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).
  - **claro** will try to get GitHub Personal Access Token from: (1) operating system keyring; (2) environment var (`GH_TOKEN`); (3) **claro**'s config file (default `$HOME/.claro.env`)

You can inform the token whenever you clone repositories, or you can save the token in the operating system keyring (recommended) or in the **claro**'s config file. Please, have a look at [Authenticating with git command line to access repositories on GitHub](#authenticating-with-git-command-line-to-access-repositories-on-github).

- To save the token (**claro** will try to save to the OS keyring first and then to the config file.)
  - `claro config token add`
- To delete the token from OS keyring
  - `claro config token del`

## Authenticating with git command line to access repositories on GitHub



> Using an HTTPS remote URL has some advantages compared with using SSH. It's easier to set up than SSH, and usually works through strict firewalls and proxies. However, it also prompts you to enter your GitHub credentials every time you pull or push a repository. ([GitHub Docs](https://docs.github.com/en/get-started/getting-started-with-git/why-is-git-always-asking-for-my-password)).

To access repositories on GitHub from the command line application over HTTPS you must authenticate with a [GitHub personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#authenticating-with-the-command-line).

 You can avoid being prompted for your password (personal access token) by configuring Git to cache your credentials for you on [git credential storage](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage). Git works with several credential helper:

  - **Unsafe and simple way** (not recommended!)
    - This method stores your password in plaintext on disk, protected only by filesystem permissions (default `$HOME/.git-credentials`)
    - on Linux, macOS or Windows
      ```bash
      git config --global credential.helper store
      ```
      - The first time git will ask the username and password. Subsequent request will use the credentials from the store (default `$HOME/.git-credentials`). 
        - Example:
          `https://username:ghp_token@github.com`
  - **Safe with a little complexity**
    - This method stores your password encrypted on disk in the operating system keyring service. 
    - You can use native operating system keyring 
      - Linux [Secret Service](https://specifications.freedesktop.org/secret-service/latest/) dbus interface, which is provided by [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring)
        - Use [Seahorse app](https://wiki.gnome.org/Apps/Seahorse) to create a default collection with `login` name
          - Open `seahorse`; go to `file->new->password keyring`; when asked for a name, use: `login` 

        ```bash
        # ------------------------------------------------#
        # on Ubuntu Linux 22.04 LTS
        # ------------------------------------------------#
        sudo apt-get install libsecret-1-0 libsecret-1-dev g++ make
    
        cd /usr/share/doc/git/contrib/credential/libsecret && sudo make
    
        git config --global credential.helper /usr/share/doc/git/contrib/credential/libsecret/git-credential-libsecret
        # ------------------------------------------------#
        ```
      -  macOS Keychain (`/usr/bin/security`)
         -  `git config --global credential.helper osxkeychain`
    - Multi-plataform (Linux, macOS or Windows)
      - You can also use a credential helper like [Git Credential Manager](https://github.com/GitCredentialManager/git-credential-manager)