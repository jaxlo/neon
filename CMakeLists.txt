# Set the minimum required version of CMake
cmake_minimum_required(VERSION 3.10)

# Set the project name and language
project(neon CXX)

# Set the C++ standard we want to use (e.g., C++17)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Add our executable and specify its source files
# CMake will create an executable named "neon"
# from the source file "src/main.cpp"
add_executable(neon src/main.cpp)

# Add a custom target to run the executable
add_custom_target(run
    COMMAND neon
    DEPENDS neon
    WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
    COMMENT "Running neon executable"
)


# Link libraries
target_link_libraries(${PROJECT_NAME} 
    SDL3::SDL3
    Vulkan::Vulkan
)
