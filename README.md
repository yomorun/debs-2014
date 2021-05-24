# DEBS 2014

This repository accompanies our blog post "DEBS Grand Challenge: Real-Time Stream Processing with YoMo." [Link](https://blog.yomo.run/posts/debs2014-grand-challenge-by-yomo).

### Getting Started

1. Follow the instructions [here](https://yomo.run/) to install YoMo. Assuming that `$GOPATH` has been set on your device, you should be able to see a directory with the name `$GOPATH/src/github.com/yomorun/yomo`.

2. Run the following command. Don't forget to replace `${YOUR_GITHUB_USERNAME}` with your actual GitHub username!

```
mkdir -p $GOPATH/src/github.com/${YOUR_GITHUB_USERNAME} && cd $_
```

3. Clone this repository by typing `git clone git@github.com:yomorun/debs-2014.git`. Then run:

```
debs-zipper % yomo wf run workflow.yaml
debs-flow % yomo run app.go -n debs-flow
debs-source % go run main.go
```

You should be able to see some exciting results! :)

### Join Our Community

- [yomorun/yomo](https://github.com/yomorun/yomo)
- [Discord](https://discord.gg/ccuttFYC)
