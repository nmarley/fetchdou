FROM ubuntu:18.04

RUN apt-get update && apt-get -y install ghostscript

#ENTRYPOINT "/usr/bin/gs"
CMD ["/bin/bash"]

# command to extract text from PDF using GhostScript:
# gs -sDEVICE=txtwrite -o output.txt input.pdf
#CMD ["-sDEVICE=txtwrite", "-o", "output.txt", "input.pdf"]
