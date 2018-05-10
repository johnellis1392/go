#!/bin/bash

# Fetch Test Data

# The MNIST Dataset is available on Yan LeCun's website here:
# http://yann.lecun.com/exdb/mnist/
#
# The URL's for the training and test datasets are as follows:
base_url="http://yann.lecun.com/exdb/mnist"
training_images="train-images-idx3-ubyte.gz"
training_labels="train-labels-idx1-ubyte.gz"
testing_images="t10k-images-idx3-ubyte.gz"
testing_labels="t10k-labels-idx1-ubyte.gz"

fetch() {
  [[ $# == 1 ]] || return 1

  dataset="${1}"
  shift

  dataset_url="${base_url}/${dataset}"
  curl -sSL "${dataset_url}" -o "${dataset}"
  return 0
}


fetch "${training_images}"
fetch "${training_labels}"
fetch "${testing_images}"
fetch "${testing_labels}"
