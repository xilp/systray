#import "Communicator.h"  

NSInputStream *inputStream;  
NSOutputStream *outputStream;  

@implementation Communicator  

- (void)open: host port:(int)port {
    CFReadStreamRef readStream;
    CFWriteStreamRef writeStream;

    CFStreamCreatePairWithSocketToHost(NULL, (__bridge CFStringRef)host, port, &readStream, &writeStream);
    if(!CFWriteStreamOpen(writeStream)) {
        NSLog(@"Error, writeStream not open");
        return;
    }

    inputStream = (__bridge NSInputStream *)readStream;  
    outputStream = (__bridge NSOutputStream *)writeStream;  

    [inputStream setDelegate:self];  
    [outputStream setDelegate:self];  
    
    [inputStream scheduleInRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];  
    [outputStream scheduleInRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];  
    
    [inputStream open];  
    [outputStream open];  

    return;
}  

- (void)close {  
    [inputStream close];
    [outputStream close];
    
    [inputStream removeFromRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];
    [outputStream removeFromRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];
    
    [inputStream setDelegate:nil];
    [outputStream setDelegate:nil];    
}  

- (void)stream:(NSStream *)stream event:(NSStreamEvent)event {  
    switch(event) {  
        case NSStreamEventHasSpaceAvailable: {  
            if(stream == outputStream) {  
                //uint8_t *buf = (uint8_t *)[@"abc" UTF8String];  
                //[outputStream write:buf maxLength:strlen((char *)buf)];  
                //NSLog(@"Sent.");
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
                    
                    NSLog(@"%@", s);                    
                }  
            }   
            break;  
        }
        default: {
//            NSLog(@"Stream is sending an Event: %lu", event);
            break;  
        }  
    }  
}  

@end  