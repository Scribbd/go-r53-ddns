# Goals
This is a simple exercise in creating a full build from Golang to docker build to helm. With the final goal for it being to run inside my personal k8s cluster. There it will either be triggered by a cronjob or part of some other automation to allow cert-manager to do its job with a working DNS.

## Evaluation:
Nothing is learned without evaluating. The following topics were learned:
- Golang
- GitHub Actions
- Helm
### Golang
My experience in programming is as follows:
- Formal education in Java
- Some projects with JavaScript
- Some extensive projects with Python

Introduced concepts with Golang:
- Pointers / Resolvers
- Anonymous scopes
- AWS go SDK v2

I mostly have worked with languages with which I didn't have to worry about pointers, garbage collection, or reassigning variables. Go introduced me to these concepts rather quickly as I was struggling to understand if I should use a `&` or a `*` in the API calls with the AWS SDK go library. Autocompletion was a life-saver here, but I wanted to know what these symbols meant.
- `&` gives the memory address  
- `*` resolves the address to a value
I know it is a bit more complex than what is stated above. The keyword I am noting down to study later is `dereferencing`.

Anonymous scopes / anonymous functions / delta functions. They remain a mystery to me. It is a clear hole in my knowledge in programming. I did experiment here with the simplest of them the anonymous scopes. This was needed as Golang is strictly typed, and once a variable has been set, it cannot be reset to a different type. Or at least, that is what I concluded from a quick google search. So instead of having a response variable for every different API call, or having a function defined that is called only once. I tried this technique and let the variable expire by scope. This isn't possible in python, and I now wish it is.

### GitHub Actions
My experience with automation platforms:
- Some extensive projects with Jenkins
- Some projects with Ansible

To figure out:
- How to test GitHub actions offline.

No new concepts were introduced to me with GitHub Actions. GitHub Actions mostly resembles a (deceptively) simplified form of Jenkins Declarative Pipelines with the command structure/notation of Ansible.

I encountered a rather interesting (expected) behaviour with git which I never encountered before: Git does not push case sensitive filename changes (in VS Code). I created a file name `DockerFile` did some work in it, and committed this new file. Realizing I made the same naming mistake, again. I simply hit `F2` in VS Code, and pushed this to GitHub at a later date.  
However, when I started creating my workflow it would throw an error: `Dockerfile not found in directory`. Which was a massive headscratcher for me as VS Code clearly had it. But GitHub did not have `Dockerfile` it still had an imposter `DockerFile` with the content of `Dockerfile`.

Long story short: I resolved this with a `git mv` on a fresh clone.

#### Helm
My experience with Helm:
- Installing some charts on my personal k8s cluster
- Nothing else

Currently I am working on a [bare metal raspberry pi k3s cluster](https://github.com/Scribbd/k8s-in-r7y-pi). And part of that project was figuring out a solution for keeping my Dynamically assigned IP by my ISP in check.

Building a Helm chart is an experience on its own. It so is easy to fall into the I-could-do-this-better-if-I-generalise-it trap. And yes, I did fall in that trap.

I encountered a lot of troubles with this one. Taking way more time than planned. This was mostly due to my personal k8s cluster

To figure out:  
`helm package` deployment package can be made so that you can add my repository to helm and just get it directly without git. I feel like this can be a new [workflow](#github-actions).