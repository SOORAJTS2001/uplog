import logging
import time

logging.basicConfig(level=logging.INFO)
while True:
    time.sleep(1)
    logging.warning("This is the warning")
    logging.error("This is the error")
    logging.info("This is the info")
