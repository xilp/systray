#import "Communicator.h"  

CFReadStreamRef readStream;  
CFWriteStreamRef writeStream;  

NSInputStream *inputStream;  
NSOutputStream *outputStream;  

@implementation Communicator  

- (void)setup {  
    NSURL *url = [NSURL URLWithString:host];  
    
    NSLog(@"Setting up connection to %@ : %i", [url absoluteString], port);  
    
    CFStreamCreatePairWithSocketToHost(kCFAllocatorDefault, (__bridge_retained CFStringRef)[url host], port, &readStream, &writeStream);  
    
    if(!CFWriteStreamOpen(writeStream)) {  
        NSLog(@"Error, writeStream not open");  
        
        return;  
    }  
    [self open];   
    
    NSLog(@"Status of outputStream: %lu", [outputStream streamStatus]);  
    
    return;  
}  

- (void)open {  
    NSLog(@"Opening streams.");  
    
    inputStream = (__bridge_transfer NSInputStream *)readStream;  
    outputStream = (__bridge_transfer NSOutputStream *)writeStream;  
        
    [inputStream setDelegate:self];  
    [outputStream setDelegate:self];  
    
    [inputStream scheduleInRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];  
    [outputStream scheduleInRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];  
    
    [inputStream open];  
    [outputStream open];  
}  

- (void)close {  
    NSLog(@"Closing streams.");  
    
    [inputStream close];  
    [outputStream close];  
    
    [inputStream removeFromRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];  
    [outputStream removeFromRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];  
    
    [inputStream setDelegate:nil];  
    [outputStream setDelegate:nil];  
    
    inputStream = nil;  
    outputStream = nil;  
}  

- (void)stream:(NSStream *)stream handleEvent:(NSStreamEvent)event {  
    NSLog(@"Stream triggered.");  
    
    switch(event) {  
        case NSStreamEventHasSpaceAvailable: {  
            if(stream == outputStream) {  
                NSLog(@"outputStream is ready.");   
            }  
            break;  
        }  
        case NSStreamEventHasBytesAvailable: {  
            if(stream == inputStream) {  
                NSLog(@"inputStream is ready.");   
                
                uint8_t buf[1024];  
                unsigned int len = 0;  
                
                len = [inputStream read:buf maxLength:1024];  
                
                if(len > 0) {  
                    NSMutableData* data=[[NSMutableData alloc] initWithLength:0];  
                    
                    [data appendBytes: (const void *)buf length:len];  
                    
                    NSString *s = [[NSString alloc] initWithData:data encoding:NSASCIIStringEncoding];  
                    
                    [self readIn:s];                    
                }  
            }   
            break;  
        }  
        default: {  
            NSLog(@"Stream is sending an Event: %lu", event);
            break;  
        }  
    }  
}  

- (void)readIn:(NSString *)s {  
    NSLog(@"Reading in the following:");  
    NSLog(@"%@", s);  
}  

- (void)writeOut:(NSString *)s {  
    uint8_t *buf = (uint8_t *)[s UTF8String];  
    
    [outputStream write:buf maxLength:strlen((char *)buf)];  
    
    NSLog(@"Writing out the following:");  
    NSLog(@"%@", s);  
}  

@end  