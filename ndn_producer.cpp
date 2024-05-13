#include <ndn-cxx/face.hpp>
#include <iostream>

class Producer {
public:
    Producer() : face_() {}

    void run() {
        ndn::Name prefix("/ndn/testApp");
        face_.setInterestFilter(prefix,
                                std::bind(&Producer::onInterest, this, _1, _2, std::ref(face_)),
                                ndn::RegisterPrefixSuccessCallback(),
                                std::bind(&Producer::onRegisterFailed, this, _1, _2));

        // Process events until the application is stopped
        face_.processEvents();
    }

private:
    void onInterest(const ndn::Name& name, const ndn::Interest& interest, ndn::Face& face) {
        std::cout << "Received Interest: " << interest << std::endl;

        auto data = std::make_shared<ndn::Data>(interest.getName());
        std::string content = "Hello, world!";
        data->setContent(std::string_view(content));

        ndn::KeyChain keyChain;
        keyChain.sign(*data);  // Sign the data packet

        face.put(*data);  // Send the data packet
    }

    void onRegisterFailed(const ndn::Name& prefix, const std::string& reason) {
        std::cerr << "Failed to register prefix \"" << prefix.toUri() << "\": " << reason << std::endl;
    }

    ndn::Face face_;
};

int main() {
    Producer producer;
    producer.run();
    return 0;
}

