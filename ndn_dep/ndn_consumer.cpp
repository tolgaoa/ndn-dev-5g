#include <ndn-cxx/face.hpp>
#include <ndn-cxx/security/validator-config.hpp>

#include <iostream>

class Consumer {
public:
	/*
    Consumer()
    : face_("nfd-service")  // Specify NFD service hostname
    {}
	*/
    
	Consumer()
    : face_()
    {}


    void run() {
        ndn::Name interestName("/ndn/testApp");
        ndn::Interest interest(interestName);
        interest.setInterestLifetime(ndn::time::seconds(10));  // Set the interest lifetime
        interest.setMustBeFresh(true);

        std::cout << "Sending Interest: " << interest << std::endl;
        face_.expressInterest(interest,
                              std::bind(&Consumer::onData, this, std::placeholders::_1, std::placeholders::_2),
                              std::bind(&Consumer::onNack, this, std::placeholders::_1, std::placeholders::_2),
                              std::bind(&Consumer::onTimeout, this, std::placeholders::_1));

        // Process events until the application is stopped
        face_.processEvents();
    }

private:
    void onData(const ndn::Interest& interest, const ndn::Data& data) {
        std::string content(reinterpret_cast<const char*>(data.getContent().value()), data.getContent().value_size());
        std::cout << "Received Data: " << content << std::endl;
    }

    void onNack(const ndn::Interest& interest, const ndn::lp::Nack& nack) {
        std::cout << "Received Nack: " << nack.getReason() << std::endl;
    }

    void onTimeout(const ndn::Interest& interest) {
        std::cerr << "Interest timed out: " << interest.getName().toUri() << std::endl;
    }

    ndn::Face face_;
};

int main() {
    Consumer consumer;
    consumer.run();
    return 0;
}


