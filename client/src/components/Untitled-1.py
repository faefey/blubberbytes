


def is_symmetrical_title(title):
    pointer1 = 0
    

    title = title.replace(' ', '').lower()
    pointer2 = len(title)-1
    while pointer1 < pointer2: 
        if (title[pointer1] == title[pointer2]):
            pointer1 +=1
            pointer2 -=1
        else: 
            return False
    return True        

  

print(is_symmetrical_title("A Santa at NASA"))
print(is_symmetrical_title("Social Media")) 