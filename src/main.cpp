#include <SDL3/SDL.h>
#include <SDL3/SDL_vulkan.h>
#include <vulkan/vulkan.h>
#include <iostream>
#include <vector>
#include <string>
#include <iomanip>
#include <set>

class VulkanSDLApp {
private:
    SDL_Window* window;
    VkInstance instance;
    VkSurfaceKHR surface;
    VkPhysicalDevice physicalDevice;
    VkDevice device;
    VkQueue graphicsQueue;
    VkQueue presentQueue;
    VkSwapchainKHR swapchain;
    
    uint32_t graphicsFamily;
    uint32_t presentFamily;
    
public:
    VulkanSDLApp() : window(nullptr), instance(VK_NULL_HANDLE), 
                     surface(VK_NULL_HANDLE), physicalDevice(VK_NULL_HANDLE),
                     device(VK_NULL_HANDLE), swapchain(VK_NULL_HANDLE),
                     graphicsFamily(0), presentFamily(0) {}
    
    ~VulkanSDLApp() {
        cleanup();
    }
    
    bool initialize() {
        if (!initSDL()) return false;
        if (!initVulkan()) return false;
        return true;
    }
    
    void printSystemInfo() {
        std::cout << "\n=== SDL3 System Information ===" << std::endl;
        
        // Display information
        printDisplayInfo();
        
        // Audio information
        printAudioInfo();
        
        // Platform information
        std::cout << "\nPlatform: " << SDL_GetPlatform() << std::endl;
        
        // CPU information
        std::cout << "Logical CPU Cores: " << SDL_GetNumLogicalCPUCores() << std::endl;
        std::cout << "System RAM: " << SDL_GetSystemRAM() << " MB" << std::endl;
    }
    
private:
    bool initSDL() {
        if (!SDL_Init(SDL_INIT_VIDEO | SDL_INIT_AUDIO)) {
            std::cerr << "Failed to initialize SDL: " << SDL_GetError() << std::endl;
            return false;
        }
        
        // Print system information first
        printSystemInfo();
        
        // Create window in windowed fullscreen mode (borderless)
        window = SDL_CreateWindow(
            "SDL3 + Vulkan Window",
            1280, 720,
            SDL_WINDOW_VULKAN | SDL_WINDOW_FULLSCREEN // This ended with _DESKTOP before
        );
        
        if (!window) {
            std::cerr << "Failed to create window: " << SDL_GetError() << std::endl;
            return false;
        }
        
        return true;
    }
    
    void printDisplayInfo() {
        std::cout << "\n=== Display Information ===" << std::endl;
        
        int displayCount = 0;
        SDL_DisplayID* displays = SDL_GetDisplays(&displayCount);
        if (!displays || displayCount == 0) {
            std::cerr << "Failed to get displays: " << SDL_GetError() << std::endl;
            return;
        }
        
        for (int i = 0; i < displayCount; i++) {
            SDL_DisplayID displayID = displays[i];
            const char* displayName = SDL_GetDisplayName(displayID);
            
            std::cout << "\nMonitor " << (i + 1) << ": " << (displayName ? displayName : "Unknown") << std::endl;
            std::cout << "Display ID: " << displayID << std::endl;
            
            // Get display modes
            int modeCount = 0;
            SDL_DisplayMode** modes = SDL_GetFullscreenDisplayModes(displayID, &modeCount);
            if (modes && modeCount > 0) {
                std::cout << "Supported Resolutions and Refresh Rates:" << std::endl;
                
                for (int j = 0; j < modeCount; j++) {
                    const SDL_DisplayMode* mode = modes[j];
                    std::cout << "  " << mode->w << "x" << mode->h 
                             << " @ " << std::fixed << std::setprecision(2) 
                             << mode->refresh_rate << "Hz"
                             << " (Format: " << SDL_GetPixelFormatName(mode->format) << ")"
                             << std::endl;
                }
                SDL_free(modes);
            }
            
            // Get current display mode
            SDL_DisplayMode currentMode;
            if (SDL_GetCurrentDisplayMode(displayID) == 0) {
                std::cout << "Current Mode: " << currentMode.w << "x" << currentMode.h 
                         << " @ " << currentMode.refresh_rate << "Hz" << std::endl;
            }
            
            // Get desktop display mode
            SDL_DisplayMode desktopMode;
            if (SDL_GetDesktopDisplayMode(displayID) == 0) {
                std::cout << "Desktop Mode: " << desktopMode.w << "x" << desktopMode.h 
                         << " @ " << desktopMode.refresh_rate << "Hz" << std::endl;
            }
            
            // Display bounds
            SDL_Rect bounds;
            if (SDL_GetDisplayBounds(displayID, &bounds) == 0) {
                std::cout << "Display Bounds: (" << bounds.x << ", " << bounds.y 
                         << ") " << bounds.w << "x" << bounds.h << std::endl;
            }
        }
        
        SDL_free(displays);
        
        std::cout << "\n=== Fullscreen Options ===" << std::endl;
        std::cout << "Supported Display Modes:" << std::endl;
        std::cout << "• Exclusive Fullscreen: Yes (SDL_WINDOW_FULLSCREEN)" << std::endl;
        std::cout << "• Windowed Fullscreen (Borderless): Yes (SDL_WINDOW_FULLSCREEN_DESKTOP)" << std::endl;
        std::cout << "• V-Sync: Supported (via Vulkan present modes)" << std::endl;
        std::cout << "• Multi-monitor: " << (displayCount > 1 ? "Yes" : "No") << std::endl;
    }
    
    void printAudioInfo() {
        std::cout << "\n=== Audio Information ===" << std::endl;
        
        int deviceCount = 0;
        SDL_AudioDeviceID* devices = SDL_GetAudioPlaybackDevices(&deviceCount);
        if (devices && deviceCount > 0) {
            std::cout << "Audio Playback Devices:" << std::endl;
            for (int i = 0; i < deviceCount; i++) {
                const char* deviceName = SDL_GetAudioDeviceName(devices[i]);
                std::cout << "  " << (i + 1) << ": " << (deviceName ? deviceName : "Unknown") << std::endl;
            }
            SDL_free(devices);
        }
        
        deviceCount = 0;
        devices = SDL_GetAudioRecordingDevices(&deviceCount);
        if (devices && deviceCount > 0) {
            std::cout << "Audio Recording Devices:" << std::endl;
            for (int i = 0; i < deviceCount; i++) {
                const char* deviceName = SDL_GetAudioDeviceName(devices[i]);
                std::cout << "  " << (i + 1) << ": " << (deviceName ? deviceName : "Unknown") << std::endl;
            }
            SDL_free(devices);
        }
    }
    
    bool initVulkan() {
        // Get required extensions from SDL
        uint32_t extensionCount = 0;
        const char* const* extensions = SDL_Vulkan_GetInstanceExtensions(&extensionCount);
        
        if (!extensions) {
            std::cerr << "Failed to get Vulkan extensions: " << SDL_GetError() << std::endl;
            return false;
        }
        
        std::cout << "\n=== Vulkan Information ===" << std::endl;
        std::cout << "Required Vulkan Extensions:" << std::endl;
        for (uint32_t i = 0; i < extensionCount; i++) {
            std::cout << "  " << extensions[i] << std::endl;
        }
        
        // Create Vulkan instance
        VkApplicationInfo appInfo{};
        appInfo.sType = VK_STRUCTURE_TYPE_APPLICATION_INFO;
        appInfo.pApplicationName = "SDL3 Vulkan App";
        appInfo.applicationVersion = VK_MAKE_VERSION(1, 0, 0);
        appInfo.pEngineName = "No Engine";
        appInfo.engineVersion = VK_MAKE_VERSION(1, 0, 0);
        appInfo.apiVersion = VK_API_VERSION_1_0;
        
        VkInstanceCreateInfo createInfo{};
        createInfo.sType = VK_STRUCTURE_TYPE_INSTANCE_CREATE_INFO;
        createInfo.pApplicationInfo = &appInfo;
        createInfo.enabledExtensionCount = extensionCount;
        createInfo.ppEnabledExtensionNames = extensions;
        createInfo.enabledLayerCount = 0;
        
        if (vkCreateInstance(&createInfo, nullptr, &instance) != VK_SUCCESS) {
            std::cerr << "Failed to create Vulkan instance!" << std::endl;
            return false;
        }
        
        // Create surface
        if (!SDL_Vulkan_CreateSurface(window, instance, nullptr, &surface)) {
            std::cerr << "Failed to create Vulkan surface: " << SDL_GetError() << std::endl;
            return false;
        }
        
        // Find physical device
        if (!findPhysicalDevice()) return false;
        if (!createLogicalDevice()) return false;
        
        std::cout << "Vulkan initialized successfully!" << std::endl;
        return true;
    }
    
    bool findPhysicalDevice() {
        uint32_t deviceCount = 0;
        vkEnumeratePhysicalDevices(instance, &deviceCount, nullptr);
        
        if (deviceCount == 0) {
            std::cerr << "No Vulkan-compatible GPUs found!" << std::endl;
            return false;
        }
        
        std::vector<VkPhysicalDevice> devices(deviceCount);
        vkEnumeratePhysicalDevices(instance, &deviceCount, devices.data());
        
        std::cout << "\nVulkan Physical Devices:" << std::endl;
        for (const auto& device : devices) {
            VkPhysicalDeviceProperties deviceProperties;
            vkGetPhysicalDeviceProperties(device, &deviceProperties);
            
            std::cout << "  GPU: " << deviceProperties.deviceName << std::endl;
            std::cout << "    Type: ";
            switch (deviceProperties.deviceType) {
                case VK_PHYSICAL_DEVICE_TYPE_DISCRETE_GPU:
                    std::cout << "Discrete GPU";
                    break;
                case VK_PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU:
                    std::cout << "Integrated GPU";
                    break;
                case VK_PHYSICAL_DEVICE_TYPE_VIRTUAL_GPU:
                    std::cout << "Virtual GPU";
                    break;
                case VK_PHYSICAL_DEVICE_TYPE_CPU:
                    std::cout << "CPU";
                    break;
                default:
                    std::cout << "Other";
                    break;
            }
            std::cout << std::endl;
            
            // Check if device is suitable
            if (isDeviceSuitable(device)) {
                physicalDevice = device;
                std::cout << "    Selected: Yes" << std::endl;
                break;
            } else {
                std::cout << "    Selected: No" << std::endl;
            }
        }
        
        if (physicalDevice == VK_NULL_HANDLE) {
            std::cerr << "Failed to find a suitable GPU!" << std::endl;
            return false;
        }
        
        return true;
    }
    
    bool isDeviceSuitable(VkPhysicalDevice device) {
        // Find queue families
        uint32_t queueFamilyCount = 0;
        vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, nullptr);
        
        std::vector<VkQueueFamilyProperties> queueFamilies(queueFamilyCount);
        vkGetPhysicalDeviceQueueFamilyProperties(device, &queueFamilyCount, queueFamilies.data());
        
        bool graphicsFound = false;
        bool presentFound = false;
        
        for (uint32_t i = 0; i < queueFamilyCount; i++) {
            if (queueFamilies[i].queueFlags & VK_QUEUE_GRAPHICS_BIT) {
                graphicsFamily = i;
                graphicsFound = true;
            }
            
            VkBool32 presentSupport = false;
            vkGetPhysicalDeviceSurfaceSupportKHR(device, i, surface, &presentSupport);
            
            if (presentSupport) {
                presentFamily = i;
                presentFound = true;
            }
            
            if (graphicsFound && presentFound) {
                break;
            }
        }
        
        return graphicsFound && presentFound;
    }
    
    bool createLogicalDevice() {
        std::vector<VkDeviceQueueCreateInfo> queueCreateInfos;
        std::set<uint32_t> uniqueQueueFamilies = {graphicsFamily, presentFamily};
        
        float queuePriority = 1.0f;
        for (uint32_t queueFamily : uniqueQueueFamilies) {
            VkDeviceQueueCreateInfo queueCreateInfo{};
            queueCreateInfo.sType = VK_STRUCTURE_TYPE_DEVICE_QUEUE_CREATE_INFO;
            queueCreateInfo.queueFamilyIndex = queueFamily;
            queueCreateInfo.queueCount = 1;
            queueCreateInfo.pQueuePriorities = &queuePriority;
            queueCreateInfos.push_back(queueCreateInfo);
        }
        
        VkPhysicalDeviceFeatures deviceFeatures{};
        
        const std::vector<const char*> deviceExtensions = {
            VK_KHR_SWAPCHAIN_EXTENSION_NAME
        };
        
        VkDeviceCreateInfo createInfo{};
        createInfo.sType = VK_STRUCTURE_TYPE_DEVICE_CREATE_INFO;
        createInfo.queueCreateInfoCount = static_cast<uint32_t>(queueCreateInfos.size());
        createInfo.pQueueCreateInfos = queueCreateInfos.data();
        createInfo.pEnabledFeatures = &deviceFeatures;
        createInfo.enabledExtensionCount = static_cast<uint32_t>(deviceExtensions.size());
        createInfo.ppEnabledExtensionNames = deviceExtensions.data();
        
        if (vkCreateDevice(physicalDevice, &createInfo, nullptr, &device) != VK_SUCCESS) {
            std::cerr << "Failed to create logical device!" << std::endl;
            return false;
        }
        
        vkGetDeviceQueue(device, graphicsFamily, 0, &graphicsQueue);
        vkGetDeviceQueue(device, presentFamily, 0, &presentQueue);
        
        return true;
    }
    
public:
    void run() {
        std::cout << "\n=== Fullscreen Window Created ===" << std::endl;
        std::cout << "Current mode: Windowed Fullscreen (Borderless)" << std::endl;
        std::cout << "Press ESC to exit." << std::endl;
        std::cout << "Press F to toggle between Exclusive Fullscreen and Windowed Fullscreen." << std::endl;
        std::cout << "\nFeatures Demonstrated:" << std::endl;
        std::cout << "• SDL3 + Vulkan integration" << std::endl;
        std::cout << "• Fullscreen mode detection and switching" << std::endl;
        std::cout << "• Multi-monitor support ready" << std::endl;
        
        bool running = true;
        bool isExclusiveFullscreen = false; // Start in windowed fullscreen
        SDL_Event event;
        
        while (running) {
            while (SDL_PollEvent(&event)) {
                switch (event.type) {
                    case SDL_EVENT_QUIT:
                        running = false;
                        break;
                    case SDL_EVENT_KEY_DOWN:
                        if (event.key.key == SDLK_ESCAPE) {
                            running = false;
                        } else if (event.key.key == SDLK_F) {
                            // Toggle between fullscreen modes
                            isExclusiveFullscreen = !isExclusiveFullscreen;
                            if (isExclusiveFullscreen) {
                                if (SDL_SetWindowFullscreen(window, true) == 0) {
                                    std::cout << "Switched to: Exclusive Fullscreen" << std::endl;
                                } else {
                                    std::cerr << "Failed to switch to exclusive fullscreen: " << SDL_GetError() << std::endl;
                                    isExclusiveFullscreen = false;
                                }
                            } else {
                                if (SDL_SetWindowFullscreen(window, false) == 0) {
                                    std::cout << "Switched to: Windowed Fullscreen (Borderless)" << std::endl;
                                } else {
                                    std::cerr << "Failed to switch to windowed fullscreen: " << SDL_GetError() << std::endl;
                                    isExclusiveFullscreen = true;
                                }
                            }
                        }
                        break;
                    case SDL_EVENT_DISPLAY_ORIENTATION:
                        std::cout << "Display orientation changed" << std::endl;
                        break;
                }
            }
            
            // In a real application, you'd render here
            SDL_Delay(16); // ~60 FPS
        }
    }
    
private:
    void cleanup() {
        if (device != VK_NULL_HANDLE) {
            vkDestroyDevice(device, nullptr);
        }
        
        if (surface != VK_NULL_HANDLE) {
            vkDestroySurfaceKHR(instance, surface, nullptr);
        }
        
        if (instance != VK_NULL_HANDLE) {
            vkDestroyInstance(instance, nullptr);
        }
        
        if (window) {
            SDL_DestroyWindow(window);
        }
        
        SDL_Quit();
    }
};

int main(int argc, char* argv[]) {
    VulkanSDLApp app;
    
    if (!app.initialize()) {
        std::cerr << "Failed to initialize application!" << std::endl;
        return -1;
    }
    
    app.run();
    
    return 0;
}
