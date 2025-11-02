# Konuşmadan Metne Dönüştürme Aracı (spt2)

Bu araç, ses dosyalarını Google Cloud Speech-to-Text API'sini kullanarak metne dönüştürür. Sonuçları JSON, SRT altyazı ve düz metin (TXT) formatlarında dışa aktarır.

## **Önemli Not**

**DİKKAT:** Bu proje şu anda yalnızca **İngilizce (en-US)** dilindeki konuşmaları ve metinleri desteklemektedir. Farklı diller için yapılandırma dosyalarında değişiklik yapılması gerekmektedir.

## Özellikler

- Ses dosyalarını metne dönüştürme
- Google Cloud Speech-to-Text API entegrasyonu
- JSON, SRT ve TXT formatlarında çıktı üretebilme
- Konuşmacı günlüğü (diarization) desteği
- Otomatik noktalama işaretleri ve kelime zaman damgaları
- Yapılandırılabilir deşifre seçenekleri

## Gereksinimler

- Go (1.x veya üstü)
- Google Cloud Platform (GCP) Hesabı
- Aktif bir Google Cloud Storage (GCS) Bucket'ı
- FFmpeg (Ses dosyası dönüştürme için sistemde yüklü olmalıdır)

## Kurulum

1.  **Projeyi klonlayın:**
    ```bash
    git clone https://github.com/kullanici/proje-adi.git
    cd proje-adi
    ```

2.  **Gerekli Go modüllerini indirin:**
    ```bash
    go mod tidy
    ```

3.  **Google Cloud kimlik bilgilerinizi ayarlayın:**
    - Bir GCP hizmet hesabı (Service Account) oluşturun ve bir JSON anahtar dosyası indirin.
    - İndirdiğiniz JSON dosyasının yolunu `configs/default.json` dosyasındaki `google_credentials_path` alanına girin veya kendi yapılandırma dosyanızı oluşturun.

## Yapılandırma

Projenin ana yapılandırması [`configs/default.json`](configs/default.json) dosyasında bulunur. Kullanmadan önce aşağıdaki alanları kendi GCP bilgilerinizle güncellemeniz gerekir:

- `google_credentials_path`: İndirdiğiniz hizmet hesabı anahtarının yolu.
- `project_id`: GCP Proje ID'niz.
- `gcs_bucket`: Geçici dosyaların yükleneceği GCS Bucket adınız.

## Kullanım

Aracı çalıştırmak için aşağıdaki komutu kullanın:

```bash
go run cmd/main.go <ses_dosyasi_yolu>
```

**Örnek:**

```bash
go run cmd/main.go speeches/APIs\ Explained\ in\ 6\ Minutes_\ \[hltLrjabkiY\].mp3
```

Farklı bir yapılandırma dosyası belirtmek için `-config` bayrağını kullanabilirsiniz:

```bash
go run cmd/main.go -config configs/config-tr.json speeches/sample_audio.mp3
```

## Çıktı

İşlem tamamlandığında, deşifre sonuçları varsayılan olarak `output/` dizinine kaydedilir. Oluşturulan dosyalar:

- **`<dosya_adi>.json`**: Tüm deşifre verilerini içeren detaylı JSON dosyası.
- **`<dosya_adi>.srt`**: Video oynatıcılar için uygun altyazı dosyası.
- **`<dosya_adi>.txt`**: Sadece deşifre edilmiş metni içeren dosya.

## Lisans

Bu proje MIT Lisansı altında lisanslanmıştır. Detaylar için [`LICENSE`](LICENSE) dosyasına bakınız.