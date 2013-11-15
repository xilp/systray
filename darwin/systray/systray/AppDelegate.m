#import "AppDelegate.h"

@implementation AppDelegate

@synthesize window = _window;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    host = @"127.0.0.1";
    port = 6333;

    statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    [statusItem setTitle:@"D"];
    [statusItem setAction:@selector(clicked:)];
    [statusItem setHighlightMode:YES];

    headReading = 4;
    bodyReading = 0;
    bodyCache = [[NSMutableData alloc] initWithLength:0];

    [self open];
}

- (IBAction)clicked:(id)sender {
    [self open];
    NSLog(@"click!");
}

- (void)open {
    [self close];

    CFReadStreamRef readStream;
    CFWriteStreamRef writeStream;
    
    CFStreamCreatePairWithSocketToHost(NULL, (__bridge CFStringRef)host, port, &readStream, &writeStream);
    if(!CFWriteStreamOpen(writeStream)) {
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

- (void)received:(NSString *)data {
    NSError *error = nil;
    id object = [NSJSONSerialization JSONObjectWithData:data options:0 error:&error];
    if (error) {
        return;
    }
    if([object isKindOfClass:[NSDictionary class]]) {
        NSDictionary *results = object;
    }
}

- (void)stream:(NSStream *)stream handleEvent:(NSStreamEvent)event {  
    switch(event) {  
        case NSStreamEventHasSpaceAvailable: {  
            if(stream != outputStream) {
                break;
            }
            //uint8_t *buf = (uint8_t *)[@"abc" UTF8String];  
            //[outputStream write:buf maxLength:strlen((char *)buf)];  
            //NSLog(@"Sent.");
            break;  
        }

        case NSStreamEventHasBytesAvailable: {  
            if (stream != inputStream) {
                break;
            }

            if (headReading != 0) {
                NSUInteger len = [inputStream read:headCache+4-headReading maxLength:4];
                headReading -= len;
                if (headReading == 0) {
                    bodyReading = *((NSUInteger*)headCache);
                    [bodyCache setLength:0];
                    NSLog(@"Head: %d", bodyReading);
                } else {
                    break;
                }
            }

            if (bodyReading != 0) {
                uint8_t buf[bodyReading];
                NSUInteger len = [inputStream read:buf maxLength:bodyReading];
                bodyReading -= len;
                [bodyCache appendBytes: (const void *)buf length:len];
                if (bodyReading == 0) {
                    headReading = 4;
                    NSString *str = [[NSString alloc] initWithData:bodyCache encoding:NSASCIIStringEncoding];
                    NSLog(@"Body: %@", str);
                } else {
                    break;
                }
            }
            break;  
        }

        default: {
            NSLog(@"Event: %lu", event);
            break;  
        }  
    }  
}  

@end
