package sparkconf

import (
	"reflect"
	"testing"
)

func TestNewSparkConf(t *testing.T) {
	conf := NewSparkConf()
	if conf == nil {
		t.Fatal("NewSparkConf() вернул nil")
	}
	if conf.props == nil {
		t.Fatal("Инициализированная карта свойств равна nil")
	}
	if len(conf.props) != 0 {
		t.Errorf("Изначальная карта свойств должна быть пустой, получено %d элементов", len(conf.props))
	}
}

func TestSparkConf_Set(t *testing.T) {
	conf := NewSparkConf()
	returnedConf := conf.Set("key1", "value1")

	// Проверка цепочки вызовов
	if returnedConf != conf {
		t.Errorf("Set() должен возвращать тот же экземпляр SparkConf")
	}

	// Проверка значения
	if val := conf.Get("key1"); val != "value1" {
		t.Errorf("После Set('key1', 'value1') Get('key1') вернул %q, ожидалось 'value1'", val)
	}

	// Проверка перезаписи значения
	conf.Set("key1", "newValue")
	if val := conf.Get("key1"); val != "newValue" {
		t.Errorf("После Set('key1', 'newValue') Get('key1') вернул %q, ожидалось 'newValue'", val)
	}
}

func TestSparkConf_Get(t *testing.T) {
	conf := NewSparkConf()
	conf.Set("key1", "value1")

	// Проверка существующего ключа
	if val := conf.Get("key1"); val != "value1" {
		t.Errorf("Get('key1') вернул %q, ожидалось 'value1'", val)
	}

	// Проверка несуществующего ключа
	if val := conf.Get("nonexistent"); val != "" {
		t.Errorf("Get() для несуществующего ключа должен возвращать пустую строку, получено %q", val)
	}
}

func TestSparkConf_GetAll(t *testing.T) {
	conf := NewSparkConf()
	conf.Set("key1", "value1")
	conf.Set("key2", "value2")

	props := conf.GetAll()

	// Проверка, что вернулись все элементы
	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	if !reflect.DeepEqual(props, expected) {
		t.Errorf("GetAll() вернул %v, ожидалось %v", props, expected)
	}

	// Проверка, что изменение возвращенной карты не влияет на оригинал
	props["key3"] = "value3"
	if conf.Contains("key3") {
		t.Error("Изменение возвращенной карты GetAll() не должно влиять на оригинальную конфигурацию")
	}
}

func TestSparkConf_Contains(t *testing.T) {
	conf := NewSparkConf()
	conf.Set("key1", "value1")

	// Проверка существующего ключа
	if !conf.Contains("key1") {
		t.Error("Contains('key1') вернул false, ожидалось true")
	}

	// Проверка несуществующего ключа
	if conf.Contains("nonexistent") {
		t.Error("Contains('nonexistent') вернул true, ожидалось false")
	}
}

func TestSparkConf_ToCommandLineArgs(t *testing.T) {
	tests := []struct {
		name     string
		setup    map[string]string
		expected []string
	}{
		{
			name:     "пустая конфигурация",
			setup:    map[string]string{},
			expected: []string{},
		},
		{
			name: "только spark.* ключи",
			setup: map[string]string{
				"spark.app.name":        "MyApp",
				"spark.executor.memory": "2g",
			},
			expected: []string{
				"--conf", "spark.app.name=MyApp",
				"--conf", "spark.executor.memory=2g",
			},
		},
		{
			name: "смешанные ключи",
			setup: map[string]string{
				"spark.app.name": "MyApp",
				"non.spark.key":  "value",
			},
			expected: []string{
				"--conf", "spark.app.name=MyApp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := NewSparkConf()
			for k, v := range tt.setup {
				conf.Set(k, v)
			}

			args := conf.ToCommandLineArgs()

			// Проверка количества аргументов
			if len(args) != len(tt.expected) {
				t.Fatalf("ToCommandLineArgs() вернул %d аргументов, ожидалось %d: %v",
					len(args), len(tt.expected), args)
			}

			// Проверка порядка аргументов не важна, т.к. они могут быть в разном порядке из-за итерации по мапе
			// Проверяем только наличие пар "--conf key=value"
			for i := 0; i < len(args); i += 2 {
				if i+1 >= len(args) {
					t.Fatalf("Нечетное количество аргументов в результате: %v", args)
				}

				if args[i] != "--conf" {
					t.Errorf("Ожидался аргумент '--conf', получено %q", args[i])
				}

				found := false
				for j := 0; j < len(tt.expected); j += 2 {
					if j+1 < len(tt.expected) && args[i+1] == tt.expected[j+1] {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("Аргумент %q не найден в ожидаемом списке", args[i+1])
				}
			}
		})
	}
}

func TestSparkConf_Remove(t *testing.T) {
	conf := NewSparkConf()
	conf.Set("key1", "value1")

	returnedConf := conf.Remove("key1")

	// Проверка цепочки вызовов
	if returnedConf != conf {
		t.Errorf("Remove() должен возвращать тот же экземпляр SparkConf")
	}

	// Проверка удаления
	if conf.Contains("key1") {
		t.Error("После Remove('key1') ключ все еще существует")
	}

	// Проверка удаления несуществующего ключа
	conf.Remove("nonexistent") // Не должно вызывать ошибок
}

func TestSparkConf_SetIfMissing(t *testing.T) {
	conf := NewSparkConf()

	// Проверка для несуществующего ключа
	returnedConf := conf.SetIfMissing("key1", "value1")

	// Проверка цепочки вызовов
	if returnedConf != conf {
		t.Errorf("SetIfMissing() должен возвращать тот же экземпляр SparkConf")
	}

	if val := conf.Get("key1"); val != "value1" {
		t.Errorf("После SetIfMissing('key1', 'value1') Get('key1') вернул %q, ожидалось 'value1'", val)
	}

	// Проверка для существующего ключа
	conf.SetIfMissing("key1", "newValue")
	if val := conf.Get("key1"); val != "value1" {
		t.Errorf("После SetIfMissing для существующего ключа значение изменилось на %q, должно остаться 'value1'", val)
	}
}

func TestSparkConf_Merge(t *testing.T) {
	conf1 := NewSparkConf()
	conf1.Set("key1", "value1")
	conf1.Set("common", "conf1Value")

	conf2 := NewSparkConf()
	conf2.Set("key2", "value2")
	conf2.Set("common", "conf2Value")

	returnedConf := conf1.Merge(conf2)

	// Проверка цепочки вызовов
	if returnedConf != conf1 {
		t.Errorf("Merge() должен возвращать тот же экземпляр SparkConf")
	}

	// Проверка значений после слияния
	expected := map[string]string{
		"key1":   "value1",
		"key2":   "value2",
		"common": "conf2Value", // Должно быть перезаписано значением из conf2
	}

	props := conf1.GetAll()
	if !reflect.DeepEqual(props, expected) {
		t.Errorf("После Merge() получено %v, ожидалось %v", props, expected)
	}

	// Проверка, что второй объект не изменился
	if val := conf2.Get("key1"); val != "" {
		t.Errorf("Второй объект должен остаться неизменным, но содержит ключ 'key1'")
	}
}

func TestSparkConf_ConcurrentAccess(t *testing.T) {
	// Этот тест просто проверяет, что параллельный доступ не вызывает паники
	// Для реального тестирования на гонки условий используйте -race детектор
	conf := NewSparkConf()
	done := make(chan bool)

	// Параллельный доступ на чтение
	go func() {
		for i := 0; i < 100; i++ {
			_ = conf.Get("key")
			_ = conf.Contains("key")
			_ = conf.GetAll()
			_ = conf.ToCommandLineArgs()
		}
		done <- true
	}()

	// Параллельный доступ на запись
	go func() {
		for i := 0; i < 100; i++ {
			conf.Set("key", "value")
			conf.Remove("key")
			conf.SetIfMissing("key", "value")
		}
		done <- true
	}()

	<-done
	<-done
}
