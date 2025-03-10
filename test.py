import requests
import json
import os
import sqlite3

class class_temp_db:
    def __enter__(self):
        self.conn = sqlite3.connect(
            os.path.join('.', 'data', 'temp.db'),
            check_same_thread = False,
            isolation_level = None
        )

        return self.conn

    def __exit__(self, exc_type, exc_value, traceback):
        self.conn.commit()
        self.conn.close()

def python_to_golang_sync(func_name, other_set = {}):
    with class_temp_db() as m_conn:
        m_curs = m_conn.cursor()
        
        if other_set == {}:
            other_set = func_name + ' {}'
        else:
            other_set = func_name + ' ' + json.dumps(other_set)
    
        m_curs.execute('select data from temp where name = "setup_golang_port"')
        db_data = m_curs.fetchall()
        db_data = db_data[0][0] if db_data else "3001"
    
        while 1:
            res = requests.post('http://localhost:' + db_data + '/', data = other_set)
            data = res.text

            if "error" == data:
                raise
            else:
                return data