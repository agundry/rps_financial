# rps_financial [![CircleCI Build Status](https://circleci.com/gh/agundry/rps_financial.svg?style=shield)](https://circleci.com/gh/agundry/rps_financial)

Before going any further, let's set some expectations: This is a dumb application. That said, it is a cloud distributed, highly available dumb application written in golang that is scalable with Kubernetes.

So what does it do?

My girlfriend and I distribute the financial burdens of life via games of rock paper scissors, henceforth referred to as RPS.

Grocery shopping? RPS. Date night? RPS.

Somehow, I almost always lose.

This is an api meant to enable tracking these games and the costs associated, hopefully one day answering the question: which of rock, paper, and scissors is the most fiscally responsible choice?

# Dev setup

Make sure docker and docker-compose are installed, then in the root directory run `docker-compose up -d` to initialize the database.
