# Quick Start Guide

Bu kısa kılavuz, Inspector Gadget Management uygulamasını K3s cluster'ınızda çalıştırmak için gereken tüm adımları içerir.

## Ön Koşullar

1. Çalışan bir K3s cluster
2. kubectl ve kubectl-gadget kurulu
3. Podman veya Docker kurulu

## Adım 1: Inspector Gadget'ı Kur

Eğer henüz kurmadıysanız:

```bash
kubectl gadget deploy
```

Kurulumu kontrol edin:

```bash
kubectl gadget version
kubectl get pods -n gadget
```

## Adım 2: Projeyi Klonlayın veya Oluşturun

Proje dizinine gidin:

```bash
cd inspector-gadget-management
```

## Adım 3: Container Image'larını Oluşturun

Build scripti otomatik olarak Podman veya Docker'ı tespit eder ve linux/amd64 platformu için build yapar:

```bash
./build.sh
```

**Not:** Apple Silicon (M1/M2/M3) Mac kullanıyorsanız, script otomatik olarak `--platform linux/amd64` kullanarak AMD64 tabanlı Kubernetes cluster'ları ile uyumlu image'lar oluşturur.

K3s için image'ları import edin:

```bash
# Podman kullanıyorsanız
podman save gadget-backend:latest | sudo k3s ctr images import -
podman save gadget-frontend:latest | sudo k3s ctr images import -

# Veya Makefile ile
make import-k3s

# K3d kullanıyorsanız
k3d image import gadget-backend:latest gadget-frontend:latest -c <cluster-name>
```

## Adım 4: Kubernetes'e Deploy Edin

```bash
./deploy.sh
```

## Adım 5: Deployment'ı Kontrol Edin

```bash
kubectl get pods -n gadget-management
```

Podlar çalışana kadar bekleyin (her ikisi de Running olmalı).

## Adım 6: Uygulamaya Erişin

### Yöntem 1: Port Forward (Önerilen)

```bash
kubectl port-forward -n gadget-management svc/frontend 3000:80
```

Tarayıcınızda açın: http://localhost:3000

### Yöntem 2: NodePort

K3s için direkt erişim:

```bash
# Tarayıcıda açın
http://localhost:30080
```

## İlk Gadget'ınızı Çalıştırın

1. **Gadget Seçin**: Açılır menüden "Trace Exec" seçin

2. **Parametreleri Girin** (opsiyonel):
   - Namespace: `default` (veya boş bırakın)
   - Pod Name: Boş bırakın (tüm podları izlemek için)

3. **Start Gadget** butonuna tıklayın

4. **Test Edin**: Yeni bir terminal açın ve test komutu çalıştırın:
   ```bash
   kubectl run test-pod --image=busybox --rm -it -- sh
   # Container içinde:
   ls
   pwd
   ps
   ```

5. **Sonuçları Görün**: UI'da exec olaylarını gerçek zamanlı olarak göreceksiniz!

## Diğer Gadget'ları Deneyin

### Trace TCP

TCP bağlantılarını izleyin:

```bash
# Test için bir pod başlatın
kubectl run curl-test --image=curlimages/curl --rm -it -- sh
# Container içinde:
curl https://www.google.com
```

### Snapshot Process

Çalışan işlemlerin anlık görüntüsünü alın:

1. UI'da "Snapshot Process" seçin
2. Start butonuna tıklayın
3. Tüm çalışan işlemleri görün

## Sorun Giderme

### Podlar başlamıyor

Logları kontrol edin:

```bash
kubectl logs -n gadget-management -l app=gadget-backend
```

### "kubectl-gadget not found" hatası

Backend pod'unun kubectl-gadget binary'sine sahip olduğunu kontrol edin:

```bash
kubectl exec -n gadget-management deploy/backend -- which kubectl-gadget
```

### Permission hatası

RBAC yapılandırmasını kontrol edin:

```bash
kubectl get clusterrole gadget-backend-role
kubectl get clusterrolebinding gadget-backend-binding
```

### Image pull hatası

Image'ların doğru import edildiğini kontrol edin:

```bash
# K3s için
sudo k3s ctr images ls | grep gadget

# Podman ile local image'ları kontrol et
podman images | grep gadget

# K3d için
docker exec <k3d-node> crictl images | grep gadget
```

## Temizlik

Uygulamayı kaldırmak için:

```bash
kubectl delete namespace gadget-management
```

## Sonraki Adımlar

- Farklı namespace'lerdeki podları filtreleyin
- Birden fazla gadget session'ı aynı anda çalıştırın
- Çıktıları analiz edin ve güvenlik olaylarını tespit edin

## Yardım

Daha fazla bilgi için [README.md](README.md) dosyasına bakın.
