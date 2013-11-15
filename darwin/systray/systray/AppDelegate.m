#import "AppDelegate.h"

@implementation AppDelegate

@synthesize window = _window;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    host = @"127.0.0.1";
    port = 6333;

    statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    [statusItem setAction:@selector(clicked:)];
    [statusItem setHighlightMode:YES];

    timer = [NSTimer scheduledTimerWithTimeInterval:3 target:self selector:@selector(handleTimer:) userInfo:nil repeats:YES];
    unconnected = YES;

    icons = [[NSMutableDictionary alloc] init];

    headReading = 4;
    bodyReading = 0;
    bodyCache = [[NSMutableData alloc] initWithLength:0];

    [self open];
}

- (IBAction)clicked:(id)sender {
    NSMutableDictionary *cmd = [NSMutableDictionary dictionaryWithObjectsAndKeys:@"clicked", @"action", nil];
    [self send:cmd];
}

- (void)send:(NSMutableDictionary*)cmd {
    NSError *error = NULL; 
    NSData *json = [NSJSONSerialization dataWithJSONObject:cmd options:0 error:&error];
    
    if (!json) {
        NSLog(@"Error: %@", error);
        return;
    }
    uint32 len = [json length];
    NSMutableData* cache = [[NSMutableData alloc] initWithLength:0];
    [cache appendBytes:&len length:4];
    [cache appendBytes:[json bytes] length:len];
    [outputStream write:[cache bytes] maxLength:len+4];  
}

- (void)open {
    if (!unconnected) {
        return;
    }
    unconnected = NO;

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
}  

- (void)close {  
    [inputStream close];
    [outputStream close];
    
    [inputStream removeFromRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];
    [outputStream removeFromRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];
    
    [inputStream setDelegate:nil];
    [outputStream setDelegate:nil];    
}  

- (void)received:(NSData*)data {
    NSError *error = nil;
    id object = [NSJSONSerialization JSONObjectWithData:data options:NSJSONReadingMutableLeaves error:&error];
    if (error) {
        return;
    }
    if (![object isKindOfClass:[NSDictionary class]]) {
        return;
    }
    NSDictionary *cmd = object;
    NSString *action = [cmd objectForKey:@"action"];
    if ([action isEqualToString:@"show"]) {
        NSString *path = [cmd objectForKey:@"path"];
        NSImage *icon = [icons objectForKey:path];
        if (!icon) {
            icon = [[NSImage alloc] initWithContentsOfFile:path];
            if (icon) {
                [icons setValue:icon forKey:path];
            }
        }
        if (icon) {
            [statusItem setImage:icon];
        }
        NSString *hint = [cmd objectForKey:@"hint"];
        if (hint) {
            [statusItem setToolTip:hint];
        }
    } else if ([action isEqualToString:@"exit"]) {
        [NSApp terminate:self];
    }
}

- (void)handleTimer:(NSTimer*)timer {
    [self open];
}

- (void)stream:(NSStream *)stream handleEvent:(NSStreamEvent)event {  
    switch(event) {
        case NSStreamEventHasSpaceAvailable: {  
            if(stream != outputStream) {
                break;
            }
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
                } else {
                    break;
                }
            }

            if (bodyReading != 0) {
                uint8_t buf[bodyReading];
                NSUInteger len = [inputStream read:buf maxLength:bodyReading];
                bodyReading -= len;
                [bodyCache appendBytes: (const void *)buf length:len];
                if (bodyReading != 0) {
                    break;
                }
                headReading = 4;
                [self received: bodyCache];
            }
            break;  
        }

        case NSStreamEventOpenCompleted: {
            [statusItem setTitle:@""];
            break;
        }
        case NSStreamEventEndEncountered:
        case NSStreamEventErrorOccurred: {
            [statusItem setTitle:@"!"];
            unconnected = YES;
            [self close];
            break;
        }
        default: {
            NSLog(@"Event: %lu", event);
            break;  
        }  
    }  
}  

@end
