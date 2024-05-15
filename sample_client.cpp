#include <ndn-cxx/face.hpp>
#include <iostream>
#include <thread>

int main() {
    try {
        ndn::Face face("127.0.0.1");

        ndn::Name name("/test/prefix");
        std::cout << "Registering prefix " << name << std::endl;

        face.setInterestFilter(name,
                               [](const ndn::Name& name, const ndn::Interest& interest) {
                                   std::cout << "Received Interest " << interest << std::endl;
                               },
                               [](const ndn::Name& prefix, const std::string& reason) {
                                   std::cerr << "Failed to register prefix " << prefix << " due to " << reason << std::endl;
                               });

        // process events until the user presses Ctrl+C
        face.processEvents(ndn::time::milliseconds::zero(), true);
    } catch (const std::exception& e) {
        std::cerr << "ERROR: " << e.what() << std::endl;
    }

    return 0;
}

