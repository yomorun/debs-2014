# DEBS 2014 Grand Challenge - YoMo Implementation

This repository accompanies our blog post "[DEBS Grand Challenge: Real-Time Stream Processing with YoMo](https://blog.yomo.run/posts/debs2014-grand-challenge-by-yomo)."

## Introduction

[DEBS](https://dl.acm.org/conference/debs) or the **ACM International Conference on Distributed Event Based Systems** aims to "provide a forum dedicated to the dissemination of original research, the discussion of practical insights, and the reporting of experiences relevant to event based computing that were previously scattered across several scientific and professional communities" ([reference](http://www.wikicfp.com/cfp/servlet/event.showcfp?eventid=34432)).

[The DEBS 2014 Grand Challenge - Smart Grid](https://debs.org/grand-challenges/2014/) is the 8th ACM International Conference on Distributed Event Based Systems. It focused on two problems which are relevant to the industry: real-time load prediction and anomaly detection. The data for the challenge was collected from a number of smart-home installations deployed in Germany.

## System Architecture

With [YoMo](https://github.com/yomorun/yomo), we can build a real-time stream processing system like the following:

![debs 2014 grand challenge by yomo streaming serverless framework](debs2014-arch.png)

## Getting Started

1. Install `yomo` CLI: `go install -v github.com/yomorun/yomo/cmd/yomo@latest && yomo --version`.
2. Clone this repository by typing `git clone git@github.com:yomorun/debs-2014.git && cd debs-2014`. 
3. Start `yomo-zipper`: `yomo wf run workflow.yaml`
4. Start streaming function [load-prediction](): `go run cmd/stream-fn-realtime-load-prediction/load-prediction.go`
5. Start mock `data-source` to generate IoT data: `go run cmd/mock-data-source/mock-data.go`
4. Start streaming function [outliers](): `go run cmd/stream-fn-realtime-outliers/outliers.go`

You should be able to see some exciting results! :)

### Join Our Community

If you have any questions, let us know! We hope to see you around.

- [yomorun/yomo](https://github.com/yomorun/yomo)
- [Discord](https://discord.gg/ccuttFYC)
