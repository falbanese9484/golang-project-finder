#!/bin/bash

go build -o findit .

chmod +x findit

mv findit ~/.project-finder/bin