FROM ubuntu

RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

RUN apt-get update

RUN apt-get install -y python3
RUN apt-get install -y libopencv-dev python3-opencv

RUN apt-get install -y python3-pip
RUN apt-get install -y cmake
RUN pip3 install dlib
RUN pip3 install scipy

RUN apt-get install -y git
RUN git clone https://github.com/wuhuikai/FaceSwap.git

RUN mv /FaceSwap/imgs/*.jpg /FaceSwap/

WORKDIR /FaceSwap

ENTRYPOINT ["python3", "main.py"]
CMD ["--src", "test6.jpg", "--dst", "test7.jpg", "--out", "output6_7.jpg", "--correct_color"]
