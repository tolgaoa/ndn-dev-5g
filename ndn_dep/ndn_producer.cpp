#include <ndn-cxx/face.hpp>
#include <iostream>
#include <thread>  // Include for std::this_thread

class Producer {
public:
    Producer() : face_("127.0.0.1") {} // Specify the K8s service for NFD

    void run() {
        ndn::Name prefix("/ndn/testApp");
        bool isRegistered = false;

        // Attempt to set interest filter and process events
        face_.setInterestFilter(prefix,
                                std::bind(&Producer::onInterest, this, std::placeholders::_1, std::placeholders::_2, std::ref(face_)),
                                [this, &isRegistered](const ndn::Name& registeredPrefix) {
                                    std::cout << "Successfully registered prefix: " << registeredPrefix << std::endl;
                                    isRegistered = true;
                                },
                                std::bind(&Producer::onRegisterFailed, this, std::placeholders::_1, std::placeholders::_2));

        // Process events until the application successfully registers the prefix
        while (!isRegistered) {
            try {
                face_.processEvents();
                std::this_thread::sleep_for(std::chrono::seconds(1));
            } catch (const std::exception& e) {
                std::cerr << "Error during processing events: " << e.what() << std::endl;
                std::cout << "Attempting to re-register prefix..." << std::endl;
                // Optionally attempt to re-register the prefix or handle other recovery here
            }
        }

        // Once registered, you might want to keep the application running or handle additional tasks
        std::cout << "Prefix registration succeeded, continuing normal operations..." << std::endl;
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
        std::cout << "Re-attempting to register prefix after a brief pause..." << std::endl;
        std::this_thread::sleep_for(std::chrono::seconds(5)); // Pause before retrying
    }

    ndn::Face face_;
};

int main() {
    Producer producer;
    producer.run();
    return 0;
}

